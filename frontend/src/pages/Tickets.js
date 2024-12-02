import React, { useState, useEffect } from "react";
import axios from "../services/api";
import DataTable from "../components/DataTable";

function Tickets() {
  const [tickets, setTickets] = useState([]);

  useEffect(() => {
    const fetchTickets = async () => {
      const response = await axios.get("/tickets");
      setTickets(response.data);
    };
    fetchTickets();
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Tickets</h1>
      <DataTable
        columns={["id", "farmerId", "status", "description", "createdAt"]}
        data={tickets}
      />
    </div>
  );
}

export default Tickets;
