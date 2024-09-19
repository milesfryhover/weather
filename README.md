
# World's Best Weather App

This application provides the current and extended weather forecast for a given address. It uses geocoding to convert the address into coordinates and fetches weather data using a weather API. The application also implements caching to reduce API calls for repeated queries.

## Features

- **Address Input**: Enter an address (either a full address or an incomplete address), and the app will attempt to convert it to latitude and longitude coordinates.
- **Weather Forecast**: Get the current weather (temperature, high, and low) and an extended forecast for the next seven days.
- **Caching**: The app caches the forecast for each address by postal code and retrieves the cached result if an address with the same postal code is entered within 30 minutes.
- **Error Handling**: If there are issues with the geocode API or fetching the weather data, the app notifies the user and prompts them to try again.

## Usage

1. Clone the repository:
   ```bash
   git clone https://github.com/mfryhover/weather.git
   ```

2. Navigate to the project directory:
   ```bash
   cd weather
   ```
3. Run the application from the command line:
   ```bash
   go run .
   ```

4. The app will prompt you to enter an address. For example:
   ```
   To exit please enter q
   Otherwise, please enter your address
   -> 3001 Esperanza Crossing, Austin, TX 78758, USA
   ```

The app will display:
1. The current temperature, high, and low for the given location.
2. An extended forecast with high and low temperatures for the upcoming days.

If the same postal code is queried within 30 minutes, the app will return the cached forecast.

To exit the app, simply enter `q`.

### Getting a Google Geocoding API Key
To use the geocoding functionality of this application, you need to obtain an API key from Google Cloud.
You can find the instructions on how to get one [here](https://developers.google.com/maps/documentation/geocoding/overview).

Once you have your API Key, you need to set it as an environment variable named `GEOCODE_API_KEY`. You can do this by adding it to your shell environment.

## Testing

Unit tests are included for key components:
- `cache_test.go`: Tests the caching mechanism.
- `forecast_test.go`: Tests the forecast retrieval logic.
- `geocode_test.go`: Tests geocoding functionality.
- `main_test.go`: Tests main functionality for getForecast

To run the tests:
```bash
go test ./...
```

## Dependencies

- Go 1.22.4 or higher

## Components

The application was designed with several distinct components:
1. **Main Program (`main.go`)**:
   - This is the entry point of the program, responsible for reading user input, coordinating the weather forecast retrieval, and handling the cache.
   - It interacts with other components like the caching and API logic to retrieve weather data and display it to the user.

2. **Cache (`cache.go`)**:
   - This component implements an in-memory cache to store weather data for previously queried addresses.
   - It has methods such as `Add` to add new entries, `Get` to retrieve cached entries, and `PurgeCache` to remove stale entries based on a timer.

3. **API (`forecast.go`, `geocode.go`)**:
   - These files handle communication with external APIs to fetch geocoding information (to convert addresses to coordinates) and weather data.
   - The program uses `api.AddressToCoordinates` to convert an address into latitude and longitude, and `api.GetForecast` to fetch weather information for those coordinates.

4. **Testing (`*_test.go`)**:
   - Each functional component (e.g., cache, forecast, geocode) has its own dedicated test files to ensure that the code behaves as expected under various scenarios.

## Scalability Considerations

This application has basic scalability considerations, particularly in caching and concurrent execution:

1. **Caching Mechanism**:
   - An in-memory cache reduces the number of API calls for repeated queries, improving resource usage efficiency.
   - However, in production, a more robust cache system (e.g., Redis or Memcached) might be necessary to handle larger scales. A naive solution like this might degrade with high request rates, allocation rates, and a growing number of live objects.

2. **Concurrency**:
   - A background goroutine purges the cache periodically, allowing the app to remain responsive while managing memory resources. This prevents memory bloat and keeps performance stable.

3. **API Limits**:
   - The current implementation assumes a low-volume usage. For higher scale (e.g., a large number of users), API rate limiting could become a bottleneck. To address this, implementing retries with exponential backoff and rate-limiting strategies would be essential to handle bursts of requests gracefully.

4. **Statelessness**:
   - The application can be modified to run as a stateless service if deployed in a cloud environment, ensuring it can scale horizontally by allowing multiple instances to serve requests without depending on a single node's cache.

## Design Considerations
1. **Target Users**:
      - The application is designed for users who want a quick and simple way to get weather information for a specific address without needing to know the latitude and longitude.
      - The CLI interface is straightforward and user-friendly, making it accessible to a wide range of users.
      - However, we do expect the person running it to have basic knowledge of setting up env variables and running a Go program
2. **Performance**:
      - The application is designed to be efficient by caching weather data and geocoding results to reduce the number of API calls.
      - The cache eviction policy ensures that the cache remains up-to-date and doesn't consume excessive memory.
3. **Design Limitations**:
      - The program is implemented as a CLI for ease of use, but it could be modified to be a web service or other implementations. The code is written to be module-agnostic of the main implementation.
      - The focus is on providing weather information quickly and efficiently without unnecessary complexity.
      - Don't over-engineer the prompt but still provide a good UX and maintain best practices
4. **Singleton Pattern for Cache**:
      - The cache uses a singleton-like approach to ensure only one instance of the cache exists throughout the program, preventing duplication and ensuring consistency.
5. **Flexible User Input**:
      - The postal code is extracted programmatically from geocoded addresses rather than provided by the user, ensuring accuracy and flexibility with user input formats.
6. **Readability Over Complexity**:
      - The application uses httptest for mocking API calls instead of libraries like gomock, prioritizing readability. Different techniques might be more appropriate depending on complexity.

## Assumptions
1. **Accuracy of Google's Geocoding API**:
   - The geocoding API will return accurate latitude and longitude coordinates for most addresses.
   - If an address is incomplete or ambiguous, the API will still return valid results based on best-effort matching.
2. **Supported Postal Codes**:
   - The application assumes that the Google Geocoding API supports all postal codes provided by the user.
   - Any invalid or unsupported addresses will be handled gracefully, and the user will be prompted to try again.
3. **Timezone Handling**:
   - The app assumes that the weather forecast data returned corresponds to the timezone of the location being queried. It does not explicitly handle timezones or daylight saving time differences.

## Possible Improvements
- Adding more robust error handling and logging to provide better feedback to users and developers.
- Implementing a more sophisticated cache eviction policy based on usage patterns or memory constraints.
- Enhancing the geocoding logic to handle incomplete addresses or international addresses more effectively.
- Adding support for additional weather data (e.g., wind speed, humidity) and more detailed forecasts.
- Adding support for multiple weather APIs to provide redundancy and improve reliability.
- Implementing a more sophisticated retry and rate-limiting strategy to handle API rate limits more effectively.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for more details.