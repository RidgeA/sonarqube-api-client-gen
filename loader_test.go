package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	fixtureVersion00     = newVersion("0.0")
	fixtureVersion100100 = newVersion("100.100")
	fixtureVersion11     = newVersion("1.1")
	fixtureVersion44     = newVersion("4.4")
	fixturePackageName   = "packageName"
	fixtureHost          = "http://localhost:9000"
)

type apiDefinitionWith func(*apiDefinition)

func createAPIDefinition(options ...apiDefinitionWith) *apiDefinition {
	a := &apiDefinition{
		PackageName: fixturePackageName,
		Version:     fixtureVersion00,
		Host:        fixtureHost,
		WebServices: make([]*webService, 0, 0),
	}

	for _, option := range options {
		option(a)
	}

	return a
}

func apiDefinitionWithWebServices(ws ...*webService) apiDefinitionWith {
	return func(d *apiDefinition) {
		d.WebServices = append(d.WebServices, ws...)
	}
}

type webServiceWith func(*webService)

func createWebService(options ...webServiceWith) *webService {
	ws := &webService{
		PackageName: fixturePackageName,
		Path:        "/api/normal",
		Since:       *fixtureVersion00,
		Actions:     []*action{},
	}

	for _, option := range options {
		option(ws)
	}

	return ws
}

func webServiceWithSince(v version) webServiceWith {
	return func(ws *webService) {
		ws.Since = v
	}
}

func webServiceWithActions(a ...*action) webServiceWith {
	return func(ws *webService) {
		ws.Actions = append(ws.Actions, a...)
	}
}

type actionWith func(*action)

func createAction(options ...actionWith) *action {
	a := &action{
		Key:    "normal",
		Since:  *fixtureVersion00,
		Params: []*param{},
	}

	for _, option := range options {
		option(a)
	}

	return a
}

func actionInternal() actionWith {
	return func(a *action) {
		a.Internal = true
	}
}

func actionSince(v version) actionWith {
	return func(a *action) {
		a.Since = v
	}
}

func actionDeprecatedSince(v version) actionWith {
	return func(a *action) {
		a.DeprecatedSince = v
	}
}

type paramWith func(*param)

func createParam(options ...paramWith) *param {
	p := &param{
		Key:                "param_key",
		Since:              *fixtureVersion00,
		Description:        "Description",
		Required:           false,
		Internal:           false,
		ExampleValue:       "example_value",
		DeprecatedSince:    *fixtureVersion00,
		PossibleValues:     nil,
		DeprecatedKey:      "",
		DeprecatedKeySince: *fixtureVersion00,
		DefaultValue:       "default_value",
		MaximumValue:       0,
		MinimumLength:      0,
		MaximumLength:      0,
		MaxValuesAllowed:   0,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func paramDeprecatedSince(v version) paramWith {
	return func(p *param) {
		p.DeprecatedSince = v
	}
}

func Test_filterDefinition(t *testing.T) {
	type args struct {
		def *apiDefinition
		f   *filter
	}
	tests := []struct {
		name string
		args args
		want *apiDefinition
	}{
		{
			name: "should remove an deprecated service if 'deprecated' param is false",
			args: args{
				def: createAPIDefinition(
					apiDefinitionWithWebServices(
						createWebService(
							webServiceWithActions(
								createAction(
									actionDeprecatedSince(*fixtureVersion11),
								),
							),
						),
						createWebService(
							webServiceWithActions(createAction()),
						),
					),
				),
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion100100,
				},
			},
			want: createAPIDefinition(
				apiDefinitionWithWebServices(
					createWebService(webServiceWithActions(createAction())),
				),
			),
		},
		{
			name: "should keep an deprecated service if 'deprecated' param is true",
			args: args{
				def: createAPIDefinition(
					apiDefinitionWithWebServices(
						createWebService(
							webServiceWithActions(
								createAction(actionDeprecatedSince(*fixtureVersion11))),
						),
						createWebService(webServiceWithActions(createAction())),
					),
				),
				f: &filter{
					deprecated: true,
					internal:   false,
					version:    fixtureVersion00,
				},
			},
			want: createAPIDefinition(
				apiDefinitionWithWebServices(
					createWebService(
						webServiceWithActions(
							createAction(actionDeprecatedSince(*fixtureVersion11))),
					),
					createWebService(webServiceWithActions(createAction())),
				),
			),
		},
		{
			name: "should remove an internal service if 'internal' param is false",
			args: args{
				def: createAPIDefinition(
					apiDefinitionWithWebServices(
						createWebService(webServiceWithActions(createAction(actionInternal()))),
						createWebService(webServiceWithActions(createAction())),
					),
				),
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion100100,
				},
			},
			want: createAPIDefinition(
				apiDefinitionWithWebServices(
					createWebService(webServiceWithActions(createAction())),
				),
			),
		},
		{
			name: "should keep an internal service if 'internal' param is true",
			args: args{
				def: createAPIDefinition(
					apiDefinitionWithWebServices(
						createWebService(webServiceWithActions(createAction(actionInternal()))),
						createWebService(webServiceWithActions(createAction())),
					),
				),
				f: &filter{
					deprecated: false,
					internal:   true,
					version:    fixtureVersion00,
				},
			},
			want: createAPIDefinition(
				apiDefinitionWithWebServices(
					createWebService(webServiceWithActions(createAction(actionInternal()))),
					createWebService(webServiceWithActions(createAction())),
				),
			),
		},
		{
			name: "should remove services with version grater than target",
			args: args{
				def: createAPIDefinition(
					apiDefinitionWithWebServices(
						createWebService(
							webServiceWithSince(*fixtureVersion44),
							webServiceWithActions(createAction()),
						),
						createWebService(
							webServiceWithActions(createAction()),
						),
					),
				),
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion11,
				},
			},
			want: createAPIDefinition(
				apiDefinitionWithWebServices(
					createWebService(
						webServiceWithActions(createAction()),
					),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterDefinition(tt.args.def, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterActions(t *testing.T) {
	type args struct {
		actions []*action
		f       *filter
	}
	tests := []struct {
		name string
		args args
		want []*action
	}{
		{
			name: "should remove deprecated actions if deprecated = false",
			args: args{
				actions: []*action{
					createAction(actionDeprecatedSince(*fixtureVersion11)),
					createAction(),
				},
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion00,
				},
			},
			want: []*action{createAction()},
		},
		{
			name: "should keep deprecated actions if deprecated = true",
			args: args{
				actions: []*action{
					createAction(actionDeprecatedSince(*fixtureVersion11)),
					createAction(),
				},
				f: &filter{
					deprecated: true,
					internal:   false,
					version:    fixtureVersion00,
				},
			},
			want: []*action{
				createAction(actionDeprecatedSince(*fixtureVersion11)),
				createAction(),
			},
		},
		{
			name: "should remove internal actions if internal = false",
			args: args{
				actions: []*action{
					createAction(actionInternal()),
					createAction(),
				},
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion00,
				},
			},
			want: []*action{createAction()},
		},
		{
			name: "should keep internal actions if internal = true",
			args: args{
				actions: []*action{
					createAction(actionInternal()),
					createAction(),
				},
				f: &filter{
					deprecated: false,
					internal:   true,
					version:    fixtureVersion00,
				},
			},
			want: []*action{
				createAction(actionInternal()),
				createAction(),
			},
		},
		{
			name: "should remove action if Since greater than target version",
			args: args{
				actions: []*action{
					createAction(
						actionSince(*fixtureVersion44),
					),
					createAction(),
				},
				f: &filter{
					deprecated: false,
					internal:   true,
					version:    fixtureVersion11,
				},
			},
			want: []*action{
				createAction(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterActions(tt.args.actions, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterParams(t *testing.T) {
	type args struct {
		params []*param
		f      *filter
	}
	tests := []struct {
		name string
		args args
		want []*param
	}{
		{
			name: "should remove deprecated param",
			args: args{
				params: []*param{createParam(), createParam(paramDeprecatedSince(*fixtureVersion11))},
				f: &filter{
					deprecated: false,
					internal:   false,
					version:    fixtureVersion11,
				},
			},
			want: []*param{createParam()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterParams(tt.args.params, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTargetVersion(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/api/server/version":
			fmt.Fprintf(w, "2.3")
		case "/bad_request/api/server/version":
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
	}))

	defer ts.Close()

	type args struct {
		client  *http.Client
		host    string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return version if passed",
			args: args{
				client:  nil,
				host:    "",
				version: "2.1",
			},
			want:    "2.1",
			wantErr: false,
		},
		{
			name: "should get version from the server if version is empty",
			args: args{
				client:  http.DefaultClient,
				host:    ts.URL,
				version: "",
			},
			want:    "2.3",
			wantErr: false,
		},
		{
			name: "should return error if response status != 200",
			args: args{
				client:  http.DefaultClient,
				host:    ts.URL + "/bad_request",
				version: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTargetVersion(tt.args.client, tt.args.host, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTargetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
