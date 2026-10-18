package handler

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

// encodedM3U contains the base64-encoded content of 1
const encodedM3U = "I0VYVE0zVSB1cmwtdHZnPSJ7e0VQR19VUkx9fSIgeC10dmctdXJsPSJ7e0VQR19VUkx9fSIKCiNFWFRJTkY6LTEgdHZnLWlkPSI2NmUwYTRmZjdlNjFlZjAwMjNiNmU5ZTgiIHR2Zy1uYW1lPSJUVkoiIHR2Zy1sb2dvPSJodHRwczovL3Jqci52b2QuaW1tZXJnby50di9yanIvdXBsb2FkZWQvMThkOTE2MTUtMmU4Zi00N2M3LTgzOTktNWZlZTdjMTI4NzVjLnBuZyIgZ3JvdXAtdGl0bGU9IjFTcG90TWVkaWEiLFRWSgpodHRwczovL3JqcnR2amxvY2FsLWlvcml2ZXItY2RuLmVuY29kZXJzLmltbWVyZ28udHYvbWFzdGVyLm0zdTgKCiNFWFRJTkY6LTEgdHZnLWlkPSI2NmUwYTQ3MmZmMTUyODAwMjM1ZTZlNzUiIHR2Zy1uYW1lPSJUVkogU3BvcnRzIE5ldHdvcmsiIHR2Zy1sb2dvPSJodHRwczovL3Jqci52b2QuaW1tZXJnby50di9yanIvdXBsb2FkZWQvZGQ2NGNkMDYtZTAxNS00NDI1LWFlMTAtOTdlYWUyNzZkNGM1LnBuZyIgZ3JvdXAtdGl0bGU9IjFTcG90TWVkaWEiLFRWSiBTcG9ydHMgTmV0d29yawpodHRwczovL3JqcnR2anNubG9jYWwtaW9yaXZlci1jZG4uZW5jb2RlcnMuaW1tZXJnby50di9tYXN0ZXIubTN1OAoKI0VYVElORjotMSB0dmctaWQ9IjY2ZTBhNzk3MTRlYzc4MDAyNWJhZTRmZCIgdHZnLW5hbWU9IkpOTiIgdHZnLWxvZ289Imh0dHBzOi8vcmpyLnZvZC5pbW1lcmdvLnR2L3Jqci91cGxvYWRlZC8zYzM0ODdhOC01ZTc3LTRjM2EtYWJiMC04YTVlY2E0NzdlODYucG5nIiBncm91cC10aXRsZT0iMVNwb3RNZWRpYSIsSk5OCmh0dHBzOi8vcmpyam5uLWlvcml2ZXItY2RuLmVuY29kZXJzLmltbWVyZ28udHYvbWFzdGVyLm0zdTgKCiNFWFRJTkY6LTEgdHZnLWlkPSI2NmUwYThjYmExN2Y4MzAwMjQxYWQ2MmMiIHR2Zy1uYW1lPSJSRVRWIiB0dmctbG9nbz0iaHR0cHM6Ly9yanIudm9kLmltbWVyZ28udHYvcmpyL3VwbG9hZGVkLzkwMDA4MzIzLTYwYzAtNDdmYy04OWFhLTA5YTg5Mjg4YmU2My5wbmciIGdyb3VwLXRpdGxlPSIxU3BvdE1lZGlhIixSRVRWCmh0dHBzOi8vcmpycmV0di1pb3JpdmVyLWNkbi5lbmNvZGVycy5pbW1lcmdvLnR2L21hc3Rlci5tM3U4CgojRVhUSU5GOi0xIHR2Zy1pZD0iNjZlMGE1ZTc3ZTYxZWYwMDIzYjZlOWVlIiB0dmctbmFtZT0iQ1ZUViIgdHZnLWxvZ289Imh0dHBzOi8vcmpyLnZvZC5pbW1lcmdvLnR2L3Jqci91cGxvYWRlZC85YmE1ZDJiMy1mZDMyLTQ1OTgtOTMyMC1hMDZlM2QzMTUwOTgucG5nIiBncm91cC10aXRsZT0iMVNwb3RNZWRpYSIsQ1ZUVgpodHRwczovL3JqcmN2dHYtaW9yaXZlci1jZG4uZW5jb2RlcnMuaW1tZXJnby50di9tYXN0ZXIubTN1OAoKI0VYVElORjotMSB0dmctaWQ9IjY3Yzg2MTNjODM0MmRmMDAyMzA2ZWEwYiIgdHZnLW5hbWU9IkRheXN0YXIiIHR2Zy1sb2dvPSJodHRwczovL3Jqci52b2QuaW1tZXJnby50di9yanIvdXBsb2FkZWQvZTdlYzc2ZTItMmZhMi00YWIzLWJkNTQtNzljMmY3Y2ZmYWMyLnBuZyIgZ3JvdXAtdGl0bGU9IjFTcG90TWVkaWEiLERheXN0YXIKaHR0cHM6Ly9saXZlLW1jbC5jZG4wMS5uZXQvd2ViL3kwOXc4bzdtci9wbGF5bGlzdC5tM3U4CgojRVhUSU5GOi0xIHR2Zy1pZD0iNjZlMGFhMDRiODMzODQwMDIzNGMyZWJkIiB0dmctbmFtZT0iRmFtZSA5NUZNIiB0dmctbG9nbz0iaHR0cHM6Ly9yanIudm9kLmltbWVyZ28udHYvcmpyL3VwbG9hZGVkL2ExY2NlZTFkLWY0YTUtNGQ2NC1iYThhLTFiMTAyOWY2OTUyMS5wbmciIGdyb3VwLXRpdGxlPSIxU3BvdE1lZGlhIixGYW1lIDk1Rk0KaHR0cHM6Ly9yanJmYW1lLWlvcml2ZXItY2RuLmVuY29kZXJzLmltbWVyZ28udHYvbWFzdGVyLm0zdTgKCiNFWFRJTkY6LTEgdHZnLWlkPSI2NmUwYTkzNDE0ZWM3ODAwMjViYWU1MGQiIHR2Zy1uYW1lPSJSYWRpbyBKYW1haWNhIiB0dmctbG9nbz0iaHR0cHM6Ly9yanIudm9kLmltbWVyZ28udHYvcmpyL3VwbG9hZGVkL2NhMDM3NGEwLTliYTUtNGNlNi04ZWZhLWI4MmMxY2NmNGRhZS5wbmciIGdyb3VwLXRpdGxlPSIxU3BvdE1lZGlhIixSYWRpbyBKYW1haWNhCmh0dHBzOi8vcmpycmotaW9yaXZlci1jZG4uZW5jb2RlcnMuaW1tZXJnby50di9tYXN0ZXIubTN1OAo="

var templateM3U string

func init() {
	raw, err := base64.StdEncoding.DecodeString(encodedM3U)
	if err != nil {
		panic("failed to decode embedded m3u: " + err.Error())
	}
	templateM3U = string(raw)
}

// M3u serves the M3U playlist with CORS and dynamic EPG resolution.
func M3u(w http.ResponseWriter, r *http.Request) {
	// CORS headers to ensure web players, IPTV Smarters, Kodi, and browser apps can fetch the playlist
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Resolve EPG URL dynamically based on host or environment variable
	scheme := "https"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if r.TLS == nil && strings.HasPrefix(r.Host, "localhost") {
		scheme = "http"
	}

	epgURL := os.Getenv("EPG_URL")
	if epgURL == "" {
		if r.Host != "" {
			epgURL = scheme + "://" + r.Host + "/api/xmltv"
		} else {
			epgURL = ""
		}
	}

	var content string
	if epgURL != "" {
		content = strings.ReplaceAll(templateM3U, "{{EPG_URL}}", epgURL)
	} else {
		// If no EPG URL can be determined, remove the EPG attributes rather than outputting invalid paths
		content = strings.ReplaceAll(templateM3U, " url-tvg=\"{{EPG_URL}}\" x-tvg-url=\"{{EPG_URL}}\"", "")
		content = strings.ReplaceAll(content, "{{EPG_URL}}", "")
	}

	w.Header().Set("Content-Type", "application/x-mpegurl; charset=utf-8")
	w.Header().Set("Content-Disposition", "inline; filename=\"1spotmedia_channels.m3u\"")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
