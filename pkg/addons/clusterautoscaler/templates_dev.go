// +build dev

package clusterautoscaler

import "net/http"

// Templates contains the cluster autoscaler's YAML templates.
var Templates http.FileSystem = http.Dir("templates")
