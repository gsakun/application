package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	typesconfig "github.com/hd-Li/types/config"
	//"github.com/rancher/norman/leader"
	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/store/proxy"
	"github.com/snowzach/rotatefilehook"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	//"github.com/rancher/norman/types"
	"github.com/hd-Li/application/controller"
	projectschema "github.com/hd-Li/types/apis/project.cattle.io/v3/schema"
	projectclient "github.com/hd-Li/types/client/project/v3"
	colorable "github.com/mattn/go-colorable"
)

var (
	kubeConfig string = "./config45"
)

func init() {
	if os.Getenv("REDIS_SERVER") == "" || os.Getenv("AUTHN_ENDPOINT") == "" || os.Getenv("AUTHN_REALM") == "" || os.Getenv("PROXYIMAGE") == "" {
		log.Fatalf("Please check env settings (%s %s %s %s)", "REDIS_SERVER", "AUTHN_ENDPOINT", "AUTHN_REALM", "PROXYIMAGE")
	}
	loglevel := os.Getenv("LOG_LEVEL")
	var logLevel log.Level
	log.Infof("loglevel env is %s", loglevel)
	if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)
		logLevel = log.DebugLevel
		log.Infof("log level is %s", loglevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
		logLevel = log.InfoLevel
		log.Infoln("log level is normal")
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	})
	log.SetOutput(colorable.NewColorableStdout())
	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	log.SetReportCaller(true)
	log.AddHook(rotateFileHook)
	//todo print log with code line
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Fatalf("Get restconfig failed: %s", err.Error())
			os.Exit(1)
		}
	}

	ctx := SigTermCancelContext(context.Background())

	userContext, err := typesconfig.NewUserOnlyContext(*restConfig)
	if err != nil {
		log.Fatalf("create userContext failed, err: %s", err.Error())
		os.Exit(1)
	}
	err = SetupApplicationCRD(ctx, userContext, *restConfig)
	if err != nil {
		log.Fatalf("create application crd failed, err: %s ", err.Error())
		os.Exit(1)
	}

	controller.Register(ctx, userContext)
	err = userContext.Start(ctx)
	if err != nil {
		panic(err)
	}
	<-ctx.Done()
	/*go leader.RunOrDie(ctx, "", "application-controller", userContext.K8sClient, func(ctx context.Context) {
		err = SetupApplicationCRD(ctx, userContext, *restConfig)
		if err != nil {
			log.Fatalf("create application crd failed, err: %s ", err.Error())
			os.Exit(1)
		}

		controller.Register(ctx, userContext)
		err = userContext.Start(ctx)
		if err != nil {
			panic(err)
		}
		<-ctx.Done()
	})
	<-ctx.Done()*/
}

// SetupApplicationCRD use for init application crd
func SetupApplicationCRD(ctx context.Context, apiContext *typesconfig.UserOnlyContext, config rest.Config) error {
	//schemas := types.Schemas{}

	applicationschema := apiContext.Schemas.Schema(&projectschema.Version, projectclient.ApplicationType)
	//schemas.AddSchema(applicationschema)

	clientGetter, err := proxy.NewClientGetterFromConfig(config)
	if err != nil {
		log.Fatalf("create clientGetter error: %s", err.Error())
		return err
	}

	factory := &crd.Factory{ClientGetter: clientGetter}
	_, err = factory.CreateCRDs(ctx, typesconfig.UserStorageContext, applicationschema)

	return err
}

// SigTermCancelContext use for kill process
func SigTermCancelContext(ctx context.Context) context.Context {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-term:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx
}
