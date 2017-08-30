#!/bin/bash

istioctl mixer rule create reviews.default.svc.cluster.local reviews.default.svc.cluster.local -f mixer-rule-additional-telemetry.yaml
