package main

const tmpl = `// Generated by nats-rpc. DO NOT EDIT.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"{{ .PkgPath }}"

	"github.com/chop-dbhi/nats-rpc/log"
	"github.com/chop-dbhi/nats-rpc/transport"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"

	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

const (
	clientType = "{{ .Pkg }}-cli"
)

var (
	buildVersion string

	jsonMarshaler = &jsonpb.Marshaler{
		EmitDefaults: true,
	}

	jsonUnmarshaler = &jsonpb.Unmarshaler{}
)

func main() {
	var (
		natsAddr     string
		printVersion bool
	)

	flag.StringVar(&natsAddr, "nats.addr", "nats://127.0.0.1:4222", "NATS address.")
	flag.BoolVar(&printVersion, "version", false, "Print version.")

	flag.Parse()

	if printVersion {
		fmt.Fprintln(os.Stdout, buildVersion)
		return
	}

	// Get method.
	args := flag.Args()

	if len(args) == 0 {
		log.Fatalf("method name required")
	}

	meth := args[0]

	// Initialize base logger.
	logger, err := log.New()
	if err != nil {
		log.Fatal(err)
	}

	logger = logger.With(
		zap.String("client.type", clientType),
		zap.String("client.version", buildVersion),
	)

	// Initialize the transport layer.
	tp, err := transport.Connect(&nats.Options{
		Url: natsAddr,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Close()

	tp.SetLogger(logger)

	inp := "{}"
	if len(args) > 1 {
		inp = args[1]
	}

	inpr := bytes.NewBufferString(inp)

	{{ $Pkg := .Pkg }}

	var rep proto.Message
	ctx := context.Background()

	switch meth { {{ $Service := .Name }}{{ range .Methods }}
	case "{{ .Name }}":
		client := {{ $Pkg }}.New{{ $Service }}Client(tp)
		var req {{ $Pkg }}.{{ .InputType | base }}
		if err := jsonUnmarshaler.Unmarshal(inpr, &req); err != nil {
			log.Fatalf("json: %s", err)
		}
		rep, err = client.{{ .Name }}(ctx, &req)
		{{ end }}

	default:
		log.Fatalf("unknown method %s", meth)
	}

	if err != nil {
		if sts, ok := status.FromError(err); ok {
			out := map[string]interface{}{
				"code": sts.Code().String(),
				"message": sts.Message(),
			}
			if err := json.NewEncoder(os.Stderr).Encode(out); err != nil {
				log.Fatalf("error encoding error: %s", err)
			}
		}
		os.Exit(1)
	}

	if err := jsonMarshaler.Marshal(os.Stdout, rep); err != nil {
		log.Fatalf("error encoding response: %s", err)
	}
	fmt.Fprint(os.Stdout, "\n")
}
`
