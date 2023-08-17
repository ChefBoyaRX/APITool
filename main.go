package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func main() {
	apiClient := &APIClient{} // Create an instance of the APIClient.
	apiData, rawJSON, err := apiClient.FetchData()
	if err != nil {
		panic(err)
	}

	fmt.Println("Raw JSON Response:")
	fmt.Println(string(rawJSON))

	// Print the API Data to inspect its contents
	fmt.Println("API Data:")
	for _, series := range apiData {
		fmt.Println(series.SeriesID)
		for _, d := range series.Data {
			fmt.Println(d.Year, d.Period, d.Value)
		}
	}

	// Ensure that apiData contains data for both series
	if len(apiData) < 2 {
		panic("Not enough data retrieved for both series.")
	}

	// Get the PeriodMap from generateLineItems function
	periodMap := generateLineItems(apiData[0].Data)

	var years []int
	for year := range periodMap {
		years = append(years, year)
	}
	sort.Ints(years)

	var sortedData []Data
	for _, year := range years {
		sortedData = append(sortedData, periodMap[year])
	}

	lineChart := lineMulti(apiData, sortedData)
	// Save the chart as an HTML file
	f, err := os.Create("line.html")
	if err != nil {
		panic(err)
	}
	page := components.NewPage()
	page.AddCharts(lineChart)
	page.Render(io.MultiWriter(f))

	// Start a local web server to serve the chart
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "line.html")
	})

	port := 8080
	addr := ":" + strconv.Itoa(port)
	println("Server running on http://localhost" + addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func lineMulti(apiData []Series, data []Data) *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "BLS API Tool",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "shine",
		}),
	)

	var xData []opts.LineData
	for _, d := range data {
		year, err := strconv.Atoi(d.Year)
		if err != nil {
			// Handle the error gracefully, skipping the data point
			continue
		}
		xData = append(xData, opts.LineData{Value: year})
	}

	s1Data := generateLineItems(apiData[0].Data)
	var s1SeriesData []opts.LineData
	for _, data := range s1Data {
		value, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			panic(err)
		}
		s1SeriesData = append(s1SeriesData, opts.LineData{Value: value})
	}

	s2Data := generateLineItems(apiData[1].Data)
	var s2SeriesData []opts.LineData
	for _, data := range s2Data {
		value, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			panic(err)
		}
		s2SeriesData = append(s2SeriesData, opts.LineData{Value: value})
	}

	line.SetXAxis(xData).
		AddSeries("CUUR0000SA0", s1SeriesData).
		AddSeries("CES0500000003", s2SeriesData)
	return line
}

func generateLineItems(data []Data) []Data {
	periodMap := make(map[int]Data)

	for _, d := range data {
		year, err := strconv.Atoi(d.Year)
		if err != nil {
			// Handle the error gracefully, skipping the data point
			continue
		}

		// Check if the year already exists in the periodMap, if not, add it
		if _, found := periodMap[year]; !found {
			// Add the data entry for the first occurrence of the year
			periodMap[year] = d
		}
	}

	// Sort the data by year before returning
	var sortedData []Data
	for _, d := range periodMap {
		sortedData = append(sortedData, d)
	}

	// Sort the data by year in ascending order
	sort.Slice(sortedData, func(i, j int) bool {
		year1, _ := strconv.Atoi(sortedData[i].Year)
		year2, _ := strconv.Atoi(sortedData[j].Year) 
		return year1 < year2
	})

	// Print the PeriodMap to inspect its contents
	fmt.Println("Period Map in generateLineItems:")
	for _, period := range sortedData {
		fmt.Println(period.Year, period.Value)
	}

	return sortedData
}
