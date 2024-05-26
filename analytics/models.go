package analytics

type PageVisit struct {
	Path      string
	PageViews uint
	Visits    uint
}

type WebVital struct {
	Good             uint
	NeedsImprovement uint
	Bad              uint
}

type WebVitalGroup struct {
	Path string
	LCP  WebVital
	INP  WebVital
	FID  WebVital
	CLS  WebVital
	TTFB WebVital
}

type AnalyticsResponse struct {
	PageVisit24h PageVisit
	PageVisit7d  PageVisit
	PageVisit30d PageVisit

	WebVital24h WebVitalGroup
	WebVital7d  WebVitalGroup
	WebVital30d WebVitalGroup
}
