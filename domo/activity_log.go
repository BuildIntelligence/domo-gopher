package domo

import (
	"fmt"
	"net/url"
)

// LogEntry describes a single event recorded in Domo's Activity/Audit Log.
type LogEntry struct {
	UserName          string `json:"userName,omitempty"`
	UserID            string `json:"userId,omitempty"`
	UserType          string `json:"userType,omitempty"`
	ActorID           int    `json:"actorId,omitempty"` //long
	ActorType         string `json:"actorType,omitempty"`
	ObjectName        string `json:"objectName,omitempty"`
	ObjectID          string `json:"objectId,omitempty"`
	ObjectType        string `json:"objectType,omitempty"`
	AdditionalComment string `json:"additionalComment,omitempty"`
	Time              string `json:"time,omitempty"`
	EventText         string `json:"eventText,omitempty"`
	Device            string `json:"device,omitempty"`
	BrowserDetails    string `json:"browserDetails,omitempty"`
	IPAddress         string `json:"ipAddress,omitempty"`
}

// AuditQueryParams contains all the query params that can be set for Domo Activity Log Queries.
type AuditQueryParams struct {
	User   string // Domo User Id
	Start  int    //long, start time epoch. In Domo's sample they're using timestamp in milliseconds. i.e. epoch * 1000
	End    int    //long, end time epoch
	Limit  int    //long, default is 50
	Offset int    //long, default is 0
}

// GetActivityLogEntries returns a list of log entries based on the query settings passed. Domo highly recommends using start and end times with limit and offset to retrieve large amounts of information. Domo Endpoint Documentation https://developer.domo.com/docs/activity-log-api-reference/activity-log
func (c *Client) GetActivityLogEntries(params AuditQueryParams) ([]*LogEntry, error) {

	query := generateAuditQueryParamsString(params)
	domoURL := fmt.Sprintf("%s/v1/audit?%s", c.baseURL, query)

	var entries []*LogEntry

	err := c.get(domoURL, &entries)
	return entries, err
}

// creates the query param(s) string for Domo's Audit log API. It'll order the params alphabetically.
func generateAuditQueryParamsString(params AuditQueryParams) string {
	q := url.Values{}
	if params.End != 0 {
		end := fmt.Sprintf("%d", params.End)
		q.Add("end", end)
	}

	if params.Limit != 0 {
		lim := fmt.Sprintf("%d", params.Limit)
		q.Add("limit", lim)
	} else {
		q.Add("limit", "50")
	}

	// Default value if param is omited in the Domo API is 0 so we can use the default value of params.Offset in the event it's not set.
	offset := fmt.Sprintf("%d", params.Offset)
	q.Add("offset", offset)

	if params.Start != 0 {
		start := fmt.Sprintf("%d", params.Start)
		q.Add("start", start)
	}

	if params.User != "" {
		q.Add("user", params.User)
	}

	query := q.Encode()
	return query
}
