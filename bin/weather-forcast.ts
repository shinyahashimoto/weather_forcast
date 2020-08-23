#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { WeatherForcastStack } from '../lib/weather-forcast-stack';

const app = new cdk.App();
new WeatherForcastStack(app, 'WeatherForcastStack');
