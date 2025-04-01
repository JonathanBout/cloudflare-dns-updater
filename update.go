package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func updateDNS(ipv4 string, ipv6 string) error {
	ctx := context.Background()
	zone := strings.TrimSpace(os.Getenv("CLOUDFLARE_ZONE"))
	record := strings.TrimSpace(os.Getenv("CLOUDFLARE_RECORDS"))

	if record == "@" {
		record = zone
	}

	updatedMessage := "Updated by cloudflare-updater at " + time.Now().UTC().Format(time.RFC822)

	zoneRes, err := cf.Zones.List(ctx, zones.ZoneListParams{Name: cloudflare.String(zone)})
	if err != nil {
		return err
	} else if len(zoneRes.Result) == 0 {
		return errors.New("zone '" + zone + "' not found")
	}

	zoneId := zoneRes.Result[0].ID

	nameFilter := dns.RecordListParamsName{}

	if record != "*" && record != "" {
		nameFilter.Exact = cloudflare.F(record)
	}

	var recordPatches []dns.BatchPatchUnionParam

	if ipv4 != "" {
		err = appendIpv4Records(ctx, ipv4, updatedMessage, &recordPatches, zoneId, nameFilter)
		if err != nil {
			return err
		}
	}

	if ipv6 != "" {
		err = appendIpv6Records(ctx, ipv6, updatedMessage, &recordPatches, zoneId, nameFilter)

		if err != nil {
			return err
		}
	}

	json.NewEncoder(os.Stdout).Encode(recordPatches)

	// update the records in a single batch
	_, err = cf.DNS.Records.Batch(ctx, dns.RecordBatchParams{
		ZoneID:  cloudflare.F(zoneId),
		Patches: cloudflare.F(recordPatches),
	})

	if err != nil {
		return errors.New("error updating records: " + err.Error())
	}

	return nil
}

func appendIpv4Records(ctx context.Context, ipv4 string, updatedMessage string, recordPatches *[]dns.BatchPatchUnionParam, zoneId string, nameFilter dns.RecordListParamsName) error {
	fmt.Println("looking for A records...")
	aRecords, err := cf.DNS.Records.List(ctx, dns.RecordListParams{
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
		ZoneID: cloudflare.F(zoneId),
		Name:   cloudflare.F(nameFilter),
	})

	if err != nil {
		return errors.New("error retrieving A records: " + err.Error())
	}

	if len(aRecords.Result) == 0 {
		fmt.Println("no A records found")
		return nil
	}

	for i := range aRecords.Result {
		fmt.Println("marking A record", aRecords.Result[i].Name, "for update")
		*recordPatches = append(*recordPatches, dns.BatchPatchAParam{
			ID: cloudflare.F(aRecords.Result[i].ID),
			ARecordParam: dns.ARecordParam{
				Content: cloudflare.F(ipv4),
				Comment: cloudflare.F(updatedMessage),
			},
		})
	}

	return nil
}

func appendIpv6Records(ctx context.Context, ipv6 string, updatedMessage string, recordPatches *[]dns.BatchPatchUnionParam, zoneId string, nameFilter dns.RecordListParamsName) error {
	fmt.Println("looking for AAAA records...")

	aaaaRecords, err := cf.DNS.Records.List(ctx, dns.RecordListParams{
		Type:   cloudflare.F(dns.RecordListParamsTypeAAAA),
		ZoneID: cloudflare.F(zoneId),
		Name:   cloudflare.F(nameFilter),
	})

	if err != nil {
		return errors.New("error retrieving AAAA records: " + err.Error())
	}

	if len(aaaaRecords.Result) == 0 {
		fmt.Println("no AAAA records found")
		return nil
	}

	for i := range aaaaRecords.Result {
		fmt.Println("marking AAAA record", aaaaRecords.Result[i].Name, "for update")
		*recordPatches = append(*recordPatches, dns.BatchPatchAAAAParam{
			ID: cloudflare.F(aaaaRecords.Result[i].ID),
			AAAARecordParam: dns.AAAARecordParam{
				Content: cloudflare.F(ipv6),
				Comment: cloudflare.F(updatedMessage),
			},
		})
	}

	return nil
}
