package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (api *ApiClient) present(fqdn string, recordTxt string, apiToken string) error {
	zones, err := api.listZones(apiToken)
	if err != nil {
		return err
	}
	zone := findMatchingZone(zones, fqdn)
	if zone == nil {
		return fmt.Errorf("no zone found for FQDN: %s", fqdn)
	}
	createdRecord, err := api.createRecords(apiToken, buildTxtRecord(fqdn, recordTxt))
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
	_, err = api.deployZone(apiToken, zone.ID)
	if err != nil {
		return fmt.Errorf("failed to deploy zone: %w", err)
	}
	log.Info().
		Interface("record", createdRecord).
		Msg("record created successfully")
	return nil
}

func (api *ApiClient) cleanup(fqdn string, recordTxt string, apiToken string) error {
	zones, err := api.listZones(apiToken)
	if err != nil {
		return err
	}
	zone := findMatchingZone(zones, fqdn)
	if zone == nil {
		return fmt.Errorf("no zone found for FQDN: %s", fqdn)
	}
	records, err := api.listRecords(apiToken, zone.ID, "")
	if err != nil {
		return err
	}
	recordContent := buildBaseTxt(fqdn, recordTxt)
	record := findMatchingRecord(records, recordContent)
	if record == nil {
		return fmt.Errorf("no record found with content: %s", recordContent)
	}
	destroyedRecord, err := api.destroyRecords(apiToken, record.ID)
	if err != nil {
		return fmt.Errorf("failed to destroy record %s: %w", recordContent, err)
	}
	_, err = api.deployZone(apiToken, zone.ID)
	if err != nil {
		return fmt.Errorf("failed to deploy zone: %w", err)
	}
	log.Info().
		Interface("record", destroyedRecord).
		Msg("record destroyed successfully")

	return nil
}

func buildTxtRecord(fqdn string, txt string) string {
	return fmt.Sprintf("%s ; %s", buildBaseTxt(fqdn, txt), buildCommentTxt())
}

func buildBaseTxt(fqdn string, txt string) string {
	return fmt.Sprintf("%s IN TXT \"%s\"", fqdn, txt)
}

func buildCommentTxt() string {
	return fmt.Sprintf("dns01 %s", time.Now().Format("2006-01-02T15:04:05-0700"))
}

func findMatchingZone(zones *ListZonesResponse, fqdn string) *Zone {
	for _, z := range zones.Zones {
		if strings.Contains(fqdn, z.ZoneName) {
			return &z
		}
	}
	return nil
}

func findMatchingRecord(records *ListRecordsResponse, record string) *Record {
	for _, r := range records.Records {
		if strings.HasPrefix(r.Content, record) {
			return &r
		}
	}
	return nil
}
