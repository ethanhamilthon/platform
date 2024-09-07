package config

import "os"

var (
	NatsUrl                    string
	ControllerStream           string
	ControllerStreamPathPrefix string
	ControllerAnswer           string
	AdapterStream              string
	AdapterStreamPathPrefix    string
	AdapterBuildImage          string
	AdapterPullGithub          string
	AdapterDetectRepo          string
	AdapterCreateDockerfile    string
)

func messageInit() {
	NatsUrl = loadEnv(os.Getenv("NATS_URL"), "nats://localhost:4222")
	ControllerStream = loadEnv(os.Getenv("CONTROLLER_STREAM"), "controller")
	ControllerAnswer = loadEnv(os.Getenv("CONTROLLER_ANSWER"), "controller.answer")
	ControllerStreamPathPrefix = loadEnv(os.Getenv("CONTROLLER_STREAM_PATH_PREFIX"), "controller.>")
	AdapterStream = loadEnv(os.Getenv("ADAPTER_STREAM"), "adapter")
	AdapterStreamPathPrefix = loadEnv(os.Getenv("ADAPTER_STREAM_PATH_PREFIX"), "adapter.>")
	AdapterBuildImage = loadEnv(os.Getenv("ADAPTER_BUILD_IMAGE"), "adapter.build.image")
	AdapterPullGithub = loadEnv(os.Getenv("ADAPTER_PULL_GITHUB"), "adapter.pull.github")
	AdapterDetectRepo = loadEnv(os.Getenv("ADAPTER_DETECT_REPO"), "adapter.detect.repo")
	AdapterCreateDockerfile = loadEnv(os.Getenv("ADAPTER_CREATE_DOCKERFILE"), "adapter.create.dockerfile")
}
