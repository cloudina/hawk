package main

import (
	"github.com/MagnusFrater/helmet"
)

func CustomHelmet() *helmet.Helmet {
	// fill Helmet with custom parameters
	helmetObj := helmet.Empty()
	helmetObj.ContentSecurityPolicy = helmet.NewContentSecurityPolicy(map[helmet.CSPDirective][]helmet.CSPSource{
		helmet.DirectiveFrameAncestors: {helmet.SourceNone},
		helmet.DirectiveDefaultSrc:     {helmet.SourceNone},
	})
	helmetObj.XContentTypeOptions = helmet.XContentTypeOptionsNoSniff
	helmetObj.XDNSPrefetchControl = helmet.XDNSPrefetchControlOn
	helmetObj.XDownloadOptions = helmet.XDownloadOptionsNoOpen
	helmetObj.ExpectCT = helmet.NewExpectCT(0, false, "")
	helmetObj.FeaturePolicy = helmet.EmptyFeaturePolicy()
	helmetObj.XFrameOptions = helmet.XFrameOptionsDeny
	helmetObj.XPermittedCrossDomainPolicies = helmet.PermittedCrossDomainPoliciesNone
	helmetObj.XPoweredBy = helmet.NewXPoweredBy(true, "")
	helmetObj.ReferrerPolicy = helmet.NewReferrerPolicy(helmet.DirectiveNoReferrer)
	helmetObj.StrictTransportSecurity = helmet.NewStrictTransportSecurity(31536000, true, false)
	helmetObj.XXSSProtection = helmet.NewXXSSProtection(true, helmet.DirectiveModeBlock, "")

	return helmetObj
}
