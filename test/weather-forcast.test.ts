import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import * as WeatherForcast from '../lib/weather-forcast-stack';

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new WeatherForcast.WeatherForcastStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});
