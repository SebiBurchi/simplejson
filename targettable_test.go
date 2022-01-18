package simplejson_test

import (
	"context"
	"fmt"
	"github.com/clambin/simplejson"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tt = simplejson.TargetTable{
		"empty": {},
		"query": {
			QueryFunc: queryFunc,
		},
		"tablequery": {
			TableQueryFunc: tableQueryFunc,
		},
		"both": {
			QueryFunc:      queryFunc,
			TableQueryFunc: tableQueryFunc,
		},
	}
)

func TestTargetTable_Targets(t *testing.T) {
	assert.Equal(t, []string{"both", "query", "tablequery"}, tt.Targets())
}

func TestTargetTable_RunQuery(t *testing.T) {
	_, err := tt.RunQuery(context.Background(), "query", &simplejson.TimeSeriesQueryArgs{})
	assert.Equal(t, "not implemented", err.Error())

	_, err = tt.RunQuery(context.Background(), "tablequery", &simplejson.TimeSeriesQueryArgs{})
	assert.Equal(t, "unknown target 'tablequery' for TimeSeries Query", err.Error())

	_, err = tt.RunQuery(context.Background(), "invalid", &simplejson.TimeSeriesQueryArgs{})
	assert.Equal(t, "unknown target 'invalid' for TimeSeries Query", err.Error())
}

func TestTargetTable_RunTableQuery(t *testing.T) {
	_, err := tt.RunTableQuery(context.Background(), "tablequery", &simplejson.TableQueryArgs{})
	assert.Equal(t, "not implemented", err.Error())

	_, err = tt.RunTableQuery(context.Background(), "query", &simplejson.TableQueryArgs{})
	assert.Equal(t, "unknown target 'query' for Table Query", err.Error())

	_, err = tt.RunTableQuery(context.Background(), "invalid", &simplejson.TableQueryArgs{})
	assert.Equal(t, "unknown target 'invalid' for Table Query", err.Error())
}

func queryFunc(_ context.Context, _ string, _ *simplejson.TimeSeriesQueryArgs) (response *simplejson.TimeSeriesResponse, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func tableQueryFunc(_ context.Context, _ string, _ *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	err = fmt.Errorf("not implemented")
	return
}
