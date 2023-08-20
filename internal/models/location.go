package models

type Location struct {
	Street     string `json:"street"`
	Area       string `json:"area"`
	Area2      string `json:"area2"`
	City       string `json:"city"`
	PostalCode string `json:"code"`
}

type GeoDecoder struct {
	Info struct {
		Statuscode int `json:"statuscode"`
		Copyright  struct {
			Text         string `json:"text"`
			ImageUrl     string `json:"imageUrl"`
			ImageAltText string `json:"imageAltText"`
		} `json:"copyright"`
		Messages []interface{} `json:"messages"`
	} `json:"info"`
	Options struct {
		MaxResults        int  `json:"maxResults"`
		IgnoreLatLngInput bool `json:"ignoreLatLngInput"`
	} `json:"options"`
	Results []struct {
		ProvidedLocation struct {
			LatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"providedLocation"`
		Locations []struct {
			Street             string `json:"street"`
			AdminArea6         string `json:"adminArea6"`
			AdminArea6Type     string `json:"adminArea6Type"`
			AdminArea5         string `json:"adminArea5"`
			AdminArea5Type     string `json:"adminArea5Type"`
			AdminArea4         string `json:"adminArea4"`
			AdminArea4Type     string `json:"adminArea4Type"`
			AdminArea3         string `json:"adminArea3"`
			AdminArea3Type     string `json:"adminArea3Type"`
			AdminArea1         string `json:"adminArea1"`
			AdminArea1Type     string `json:"adminArea1Type"`
			PostalCode         string `json:"postalCode"`
			GeocodeQualityCode string `json:"geocodeQualityCode"`
			GeocodeQuality     string `json:"geocodeQuality"`
			DragPoint          bool   `json:"dragPoint"`
			SideOfStreet       string `json:"sideOfStreet"`
			LinkId             string `json:"linkId"`
			UnknownInput       string `json:"unknownInput"`
			Type               string `json:"type"`
			LatLng             struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
			DisplayLatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"displayLatLng"`
			MapUrl string `json:"mapUrl"`
		} `json:"locations"`
	} `json:"results"`
}
