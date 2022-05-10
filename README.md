It's a part of [goCovid19](https://github.com/tty8747/goCovid19)

## goCovid19

[![pipeline status](https://gitlab.ubukubu.ru/goto/goCovid19-app/badges/develop/pipeline.svg)](https://gitlab.ubukubu.ru/goto/goCovid19-app/-/commits/develop)

Plan:

- Get all data from current year using API https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/
- Store it into DB:
  - date_value
  - country_code
  - confirmed
  - deaths
  - stringency_actual
  - stringency
- Output data by country_code in form of a table and sort them by deaths (or date_value which is the same) in ascending order
