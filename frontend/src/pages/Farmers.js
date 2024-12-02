import React, { useState, useEffect } from "react";
import axios from "../services/api";
import DataTable from "../components/DataTable";

function Farmers() {
  const [farmers, setFarmers] = useState([]);

  useEffect(() => {
    const fetchFarmers = async () => {
      const response = await axios.get("/farmers");
      setFarmers(response.data);
    };
    fetchFarmers();
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Farmers</h1>
      <DataTable columns={["id", "name", "contact", "crop"]} data={farmers} />
    </div>
  );
}

export default Farmers;
