package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
)

func updateDNS(ipv4 string, ipv6 string) error {
	ctx := context.Background()
	zone := os.Getenv("CLOUDFLARE_ZONE")
	updatedMessage := cloudflare.F("Updated by cloudflare-updater at " + time.Now().UTC().Format(time.RFC822))

	zoneRes, err := cf.Zones.List(ctx, zones.ZoneListParams{Name: cloudflare.String(zone)})
	if err != nil {
		return err
	} else if len(zoneRes.Result) == 0 {
		return errors.New("zone '" + zone + "' not found")
	}

	zoneId := zoneRes.Result[0].ID

	aRecords, err := cf.DNS.Records.List(ctx, dns.RecordListParams{
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
		ZoneID: cloudflare.F(zoneId),
	})

	if err != nil {
		return errors.New("error retrieving A record: " + err.Error())
	}

	aaaaRecords, err := cf.DNS.Records.List(ctx, dns.RecordListParams{
		Type:   cloudflare.F(dns.RecordListParamsTypeAAAA),
		ZoneID: cloudflare.F(zoneId),
	})

	if err != nil {
		return errors.New("error retrieving AAAA record: " + err.Error())
	}

	var recordPatches []dns.BatchPatchUnionParam

	for i := range aRecords.Result {
		recordPatches = append(recordPatches, dns.BatchPatchARecordParam{
			ID: cloudflare.F(aRecords.Result[i].ID),
			ARecordParam: dns.ARecordParam{
				Content: cloudflare.F(ipv4),
				Comment: updatedMessage,
			},
		})
	}

	for i := range aaaaRecords.Result {
		recordPatches = append(recordPatches, dns.BatchPatchAAAARecordParam{
			ID: cloudflare.F(aaaaRecords.Result[i].ID),
			AAAARecordParam: dns.AAAARecordParam{
				Content: cloudflare.F(ipv6),
				Comment: updatedMessage,
			},
		})
	}

	_, err = cf.DNS.Records.Batch(ctx, dns.RecordBatchParams{
		ZoneID:  cloudflare.F(zoneId),
		Patches: cloudflare.F(recordPatches),
	})

	if err != nil {
		return errors.New("error updating AAAA record: " + err.Error())
	}

	return nil
}
