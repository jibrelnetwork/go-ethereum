builder(
        buildTasks: [
                //[
                //        name: "Linters",
                //        type: "lint",
                //        method: "inside",
                //        runAsUser: "root",
                //        entrypoint: "",
                //        buildStage: "builder",
                //        command: [
                //                "apk add --no-cache git",
                //                "cd /go-ethereum",
                //                "make lint"
                //        ],
                //],
                [
                        name: "Tests",
                        type: "test",
                        method: "inside",
                        runAsUser: "root",
                        entrypoint: "",
                        buildStage: "builder",
                        command: [
                                "cd /go-ethereum",
                                "make test",
                        ],
                ],
        ],
)
