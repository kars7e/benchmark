package benchmarkreporter

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"
)

type DispatchReporter struct {
	file         string
	measurements map[string]types.SpecMeasurement
	plot         bool
	plotFunc     func(measurements map[string]types.SpecMeasurement)
}

func NewDispatchReporter(filename string, shouldPlot bool, plotter func(measurements map[string]types.SpecMeasurement)) *DispatchReporter {
	records := make(map[string]types.SpecMeasurement)
	return &DispatchReporter{
		file:         filename,
		measurements: records,
		plot:         shouldPlot,
		plotFunc:     plotter,
	}
}

func (reporter *DispatchReporter) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
}

func (reporter *DispatchReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {}

func (reporter *DispatchReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {}

func (reporter *DispatchReporter) SpecWillRun(specSummary *types.SpecSummary) {}

func (reporter *DispatchReporter) SpecDidComplete(spec *types.SpecSummary) {
	if spec.IsMeasurement {
		for _, measurement := range spec.Measurements {
			reporter.measurements[spec.ComponentTexts[len(spec.ComponentTexts)-1]] = *measurement
		}
	}

}

func (reporter *DispatchReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	fmt.Printf("Collected %v measurements\n", len(reporter.measurements))
	fmt.Printf("Outputting to %v as csv...\n", reporter.file)
	var output *os.File
	defer output.Close()
	if _, err := os.Stat(reporter.file); os.IsNotExist(err) {
		fmt.Printf("Doesn't exist: %v\n", os.IsNotExist(err))
		fmt.Println(reporter.file)
		if output, err = os.Create(reporter.file); err != nil {
			fmt.Printf("Unable to create new file, %v\n", err)
			log.Fatal("Unable to create new file")
		}
	} else {
		if output, err = os.Open(reporter.file); err != nil {
			fmt.Printf("Unable to open file, %v\n", err)
			log.Fatal("Unable to open file")
		}
	}
	writer := csv.NewWriter(output)
	defer writer.Flush()
	headers := []string{"Measurement (units)", "Smallest", "Largest", "Average", "StdDeviation", "Precision", "Number of Samples Collected"}
	_ = writer.Write(headers)
	for field, measurement := range reporter.measurements {
		data := []string{
			fmt.Sprintf("%v (%v)", field, measurement.Units),
			fmt.Sprintf("%v", measurement.Smallest),
			fmt.Sprintf("%v", measurement.Largest),
			fmt.Sprintf("%v", measurement.Average),
			fmt.Sprintf("%v", measurement.StdDeviation),
			fmt.Sprintf("%v", measurement.Precision),
			fmt.Sprintf("%v", len(measurement.Results)),
		}
		if err := writer.Write(data); err != nil {
			fmt.Println("Unable to write to csv file")
			log.Fatal("Unable to write to csv file")
		}
	}
	fmt.Println("Finished outputting results to csv")
	if reporter.plot {
		reporter.plotFunc(reporter.measurements)
	}
}
