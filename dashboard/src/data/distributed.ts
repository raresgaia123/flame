export default {
    schemas: [
        {
            "name": "A simple schema for distributed training with MQTT backend",
            "description": "This implementation is on Keras using MNIST dataset.",
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
                }
            ],
            "channels": [
                {
                    "description": "Model update is sent from a trainer to another trainer",
                    "groupBy": {
                        "type": "tag",
                        "value": [
                            "us"
                        ]
                    },
                    "name": "param-channel",
                    "pair": [
                        "trainer",
                        "trainer"
                    ],
                    "funcTags": {
                        "trainer": [
                            "ring_allreduce"
                        ]
                    }
                }
            ]
        }
    ]
}
