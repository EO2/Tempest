# Tempest

## Client for Weather API

*Get Temperature & Relative Humidity ++*

Created to get data for a thermostat app, mostly for learning and experimenting.
Runs a for-ever loop and outputs json file, temperature celcius and relative humidity.
Feel free to comment :)

> ### Data from The Norwegian Meteorological Institute
> - https://api.met.no/weatherapi/
> - https://api.met.no/weatherapi/locationforecast/2.0/documentation
>
>> Norwegian Licence for Open Government Data (NLOD) 2.0 and Creative Commons 4.0 BY International licences.
>> - https://data.norge.no/nlod/en/2.0/
>> - http://creativecommons.org/licenses/by/4.0

> ### Auto JSON > Go Struct
> - http://json2struct.mervine.net/

### TODO:
- Handle 203 - Deprecated: Add warning (currently dies)
- Handle 429 - Throttling: Limit traffic / request frequency (currently dies)
- Randomize interval? (Expires header + for loop already adds some time variation)
- Reduce timeout (currently transport default)
- Remove json file in/out or make optional
- Status lookup if errors
- Storage: bson, mongo ? time series db ? graphana (historical data available..)
- Compare RH at different temperatures (for other proj)
- Sometimes gives 304, even with Expires check. Is Expires not always precise?

___

Tempest v.0.1<br>
Get Weather Forecasts from MET.no API<br>
MET Forecast: Tue, 29 Mar 2022 02:11:22 CEST >  2.5 C  84.7 rH<br>
MET Forecast: Tue, 29 Mar 2022 02:43:22 CEST >  2.5 C  84.7 rH<br>
MET: 304 - Resource Not Modified<br>
MET Forecast: Tue, 29 Mar 2022 03:13:52 CEST >  2.5 C  84.7 rH<br>
MET Forecast: Tue, 29 Mar 2022 03:14:22 CEST >  2.1 C  91.7 rH<br>