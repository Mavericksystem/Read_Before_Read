package nim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	endpoint = "https://integrate.api.nvidia.com/v1/chat/completions",
	model	 = "nvidia/nemotron-3-ultra-550b-a55b"
)