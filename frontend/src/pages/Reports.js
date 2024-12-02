import React, { useState } from "react";
import axios from "../services/api";
import SelectInput from "../components/SelectInput";

function Reports() {
  const [reportType, setReportType] = useState("daily");
  const [report, setReport] = useState("");

  const generateReport = async () => {
    const response = await axios.get(`/reports?type=${reportType}`);
    setReport(response.data);
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Reports</h1>
      <SelectInput
        label="Report Type"
        options={["daily", "weekly", "monthly", "yearly"]}
        value={reportType}
        onChange={(e) => setReportType(e.target.value)}
      />
      <button
        onClick={generateReport}
        className="bg-blue-600 text-white py-2 px-4 rounded"
      >
        Generate Report
      </button>
      {report && (
        <div className="mt-4 p-4 border border-gray-300">
          <h2 className="text-xl font-bold">Report</h2>
          <pre>{report}</pre>
        </div>
      )}
    </div>
  );
}

export default Reports;
