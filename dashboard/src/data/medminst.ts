export default {
    schemas: [
        {
            "name": "Benchmark of FedOPT Aggregators/Optimizers using MedMNIST example schema v1.0.0 via PyTorch",
            "description": "A simple example of MedMNIST using PyTorch to test out different aggregator algorithms.",
            "roles": [
                {
                    "name": "trainer",
                    "description": "It consumes the data and trains local model",
                    "isDataConsumer": true,
                    "groupAssociation": [
                        {
                            "param-channel": "us"
                        }
                    ]
                },
                {
                    "name": "aggregator",
                    "description": "It aggregates the updates from trainers",
                    "replica": 1,
                    "groupAssociation": [
                        {
                            "param-channel": "us"
                        }
                    ]
                }
            ],
            "channels": [
                {
                    "name": "param-channel",
                    "description": "Model update is sent from trainer to aggregator and vice-versa",
                    "pair": [
                        "trainer",
                        "aggregator"
                    ],
                    "groupBy": {
                        "type": "tag",
                        "value": [
                            "us"
                        ]
                    },
                    "funcTags": {
                        "trainer": [
                            "fetch",
                            "upload"
                        ],
                        "aggregator": [
                            "distribute",
                            "aggregate"
                        ]
                    }
                }
            ]
        }
    ]
}
