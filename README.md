# GitHub Repo Details

The library to fetch repo details to simplify due diligence when assessing it for further usage.

## Features

- v0.0.1:

    - Source: [GraphQL](https://docs.github.com/en/graphql/overview/explorer)

        - Repo URL

        - License type

        - Is forked

        - Is archived

        - Is disabled

        - Creation Date

        - Last Update Date

        - Adoption:

            - Stargazes:

                - Total count

                - Dynamics: count in time

            - Forks:

                - Count

                - Dynamics: count in time

            - Watchers:

                - Total count

                - Dynamics: count in time on a daily, weekly, monthly basis

        - Issues:

            - Open issues count

            - Closed issues count

            - Dynamics:

                - Count of issues in time

                - Median duration of issue resolution

        - Releases:

            - Number of releases

            - Last release date

            - Dynamics: time between release

    - Source: RestAPI

        - [Community Profile Metrics](https://docs.github.com/en/rest/reference/repository-metrics#get-community-profile-metrics)

        - [Contributors]():

            - Total count

            - Dynamics:

                - Count of contributors in time

## License

The library and tooling is distributed under [the MIT license](./LICENSE).

## Applications

### Github Trending Repos

The app generates the list of trending repos on daily and weekly basis. It has integration to [telegram](https://telegram.org/).

## Logic

1. `GET https://github.com/trending/{{Language}}?since=daily`

2. Extract the info according to the list of [features](#Features)

3. Generate the output to be displayed as markdown table

4. Publish to using telegram bot API
