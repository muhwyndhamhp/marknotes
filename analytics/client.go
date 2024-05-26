package analytics

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
)

type Client struct {
	gql         *graphql.Client
	accountTag  string
	siteTag     string
	authHeader  string
	emailHeader string
}

func NewClient(accountTag, siteTag, secretToken, email string) *Client {
	gql := graphql.NewClient("https://api.cloudflare.com/client/v4/graphql")

	return &Client{
		gql:         gql,
		accountTag:  accountTag,
		siteTag:     siteTag,
		authHeader:  "Bearer " + secretToken,
		emailHeader: email,
	}
}

func (c *Client) GetAnalytics(ctx context.Context, slug string) (*AnalyticsResponse, error) {
	req := c.getReq(slug)

	var resp map[string]interface{}

	if err := c.gql.Run(ctx, req, &resp); err != nil {
		return nil, errs.Wrap(err)
	}

	v, ok := resp["viewer"].(map[string]interface{})
	if !ok {
		return nil, errors.New("error viewer")
	}
	a, ok := v["accounts"].([]interface{})
	if !ok {
		return nil, errors.New("error accounts")
	}

	amap, ok := a[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("error first account to map")
	}

	pageVisits, err := c.parsePageVisits(amap)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	webVitals, err := c.parseWebVitals(amap)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return &AnalyticsResponse{
		PageVisit24h: pageVisits["24h"],
		PageVisit7d:  pageVisits["7d"],
		PageVisit30d: pageVisits["30d"],
		WebVital24h:  webVitals["24h"],
		WebVital7d:   webVitals["7d"],
		WebVital30d:  webVitals["30d"],
	}, nil
}

func (c *Client) parseWebVitals(resp map[string]interface{}) (map[string]WebVitalGroup, error) {
	res := map[string]WebVitalGroup{}

	d1 := c.toWebVital("webVitals24h", resp)
	d7 := c.toWebVital("webVitals7d", resp)
	d30 := c.toWebVital("webVitals30d", resp)

	res["24h"] = d1
	res["7d"] = d7
	res["30d"] = d30

	return res, nil
}

func (c *Client) toWebVital(key string, data map[string]interface{}) WebVitalGroup {
	wvl, ok := data[key].([]interface{})
	if !ok {
		return WebVitalGroup{}
	}

	if len(wvl) == 0 {
		return WebVitalGroup{}
	}

	wv, ok := wvl[0].(map[string]interface{})
	if !ok {
		return WebVitalGroup{}
	}

	attr, ok := wv["attributes"].(map[string]interface{})
	if !ok {
		return WebVitalGroup{}
	}

	d, ok := wv["data"].(map[string]interface{})
	if !ok {
		return WebVitalGroup{}
	}

	path := attr["path"].(string)

	return WebVitalGroup{
		Path: path,
		LCP: WebVital{
			Good:             parseUint(d["lcpGood"]),
			NeedsImprovement: parseUint(d["lcpNeedsImprovement"]),
			Bad:              parseUint(d["lcpPoor"]),
		},
		INP: WebVital{
			Good:             parseUint(d["inpGood"]),
			NeedsImprovement: parseUint(d["inpNeedsImprovement"]),
			Bad:              parseUint(d["inpPoor"]),
		},
		FID: WebVital{
			Good:             parseUint(d["fidGood"]),
			NeedsImprovement: parseUint(d["fidNeedsImprovement"]),
			Bad:              parseUint(d["fidPoor"]),
		},
		CLS: WebVital{
			Good:             parseUint(d["clsGood"]),
			NeedsImprovement: parseUint(d["clsNeedsImprovement"]),
			Bad:              parseUint(d["clsPoor"]),
		},
		TTFB: WebVital{
			Good:             parseUint(d["ttfbGood"]),
			NeedsImprovement: parseUint(d["ttfbNeedsImprovement"]),
			Bad:              parseUint(d["ttfbPoor"]),
		},
	}
}

func parseUint(i interface{}) uint {
	if i == nil {
		return 0
	}

	d, ok := i.(float64)
	if !ok {
		return 0
	}

	return uint(d)
}

func (c *Client) parsePageVisits(resp map[string]interface{}) (map[string]PageVisit, error) {
	res := map[string]PageVisit{}

	d1 := c.toPageVisit("pageVisits24h", resp)

	d7 := c.toPageVisit("pageVisits7d", resp)

	d30 := c.toPageVisit("pageVisits30d", resp)

	res["24h"] = d1
	res["7d"] = d7
	res["30d"] = d30

	return res, nil
}

func (c *Client) toPageVisit(key string, data map[string]interface{}) PageVisit {
	odm, ok := data[key].([]interface{})
	if !ok {
		return PageVisit{}
	}

	if len(odm) == 0 {
		return PageVisit{}
	}

	od, ok := odm[0].(map[string]interface{})
	if !ok {
		return PageVisit{}
	}

	attr, ok := od["attributes"].(map[string]interface{})
	if !ok {
		return PageVisit{}
	}

	visits, ok := od["visits"].(map[string]interface{})
	if !ok {
		return PageVisit{}
	}

	path := attr["path"].(string)
	pageView := uint(od["pageViews"].(float64))
	visit := uint(visits["visits"].(float64))

	return PageVisit{
		Path:      path,
		PageViews: pageView,
		Visits:    visit,
	}
}

const queryBegin = `
      query {
        viewer {
          accounts(filter: { accountTag: "%s" }) {
`

const queryEnd = `
          }
        }
      }
`

func (c *Client) getReq(slug string) *graphql.Request {
	query := fmt.Sprintf(queryBegin, c.accountTag)

	pageVisitKeys := []string{"pageVisits24h", "pageVisits7d", "pageVisits30d"}
	webVitalKeys := []string{"webVitals24h", "webVitals7d", "webVitals30d"}
	startValues := []time.Time{time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, -1, 0)}

	path := fmt.Sprintf("/articles/%s.html", slug)

	fmt.Println(path)

	for i := range pageVisitKeys {
		query += c.pageVisitQuery(
			startValues[i],
			time.Now(),
			pageVisitKeys[i],
			c.siteTag,
			path,
		)
	}

	for i := range webVitalKeys {
		query += c.webVitalsQuery(
			startValues[i],
			time.Now(),
			webVitalKeys[i],
			c.siteTag,
			path,
		)
	}

	query += queryEnd

	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("X-AUTH-EMAIL", c.emailHeader)

	return req
}

const webVitalQuery = `
      %s: rumWebVitalsEventsAdaptiveGroups(
        filter: {
          AND: [
            {
              datetime_geq: "%s"
              datetime_leq: "%s"
            }
            { OR: [{ siteTag: "%s" }] }
            { bot: 0 }
            { requestPath: "%s" }
          ]
        }
        limit: 1
      ) {
        data: sum {
          lcpGood
          lcpNeedsImprovement
          lcpPoor
          ttfbGood
          ttfbNeedsImprovement
          ttfbPoor
          clsGood
          clsNeedsImprovement
          clsPoor
          fidGood
          fidNeedsImprovement
          fidPoor
          inpGood
          inpNeedsImprovement
          inpPoor
        }
        attributes: dimensions { path: requestPath }
      }
`

func (c *Client) webVitalsQuery(start, end time.Time, key, siteTag, path string) string {
	return fmt.Sprintf(
		webVitalQuery,
		key,
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
		siteTag,
		path,
	)
}

const pageVisitQuery = `
      %s: rumPageloadEventsAdaptiveGroups(
        filter: {
          AND: [
            {
              datetime_geq: "%s"
              datetime_leq: "%s"
            }
            { OR: [{ siteTag: "%s" }] }
            { bot: 0 }
            { requestPath: "%s" }
          ]
        }
        limit: 1
      ) {
        pageViews: count
        visits: sum { visits }
        attributes: dimensions { path: requestPath }
      }
`

func (c *Client) pageVisitQuery(start, end time.Time, key, siteTag, path string) string {
	return fmt.Sprintf(
		pageVisitQuery,
		key,
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
		siteTag,
		path,
	)
}
