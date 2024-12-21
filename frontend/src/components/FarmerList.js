import React, { useState, useEffect } from 'react';
import { farmerService } from '../services/farmerService';

const FarmersList = () => {
  const [farmers, setFarmers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadFarmers();
  }, []);

  const loadFarmers = async () => {
    try {
      setLoading(true);
      const data = await farmerService.fetchFarmers();
      setFarmers(data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch farmers. Please try again later.');
      console.error('Error:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h2>Farmers List</h2>
      <ul>
        {farmers.map((farmer) => (
          <li key={farmer.id}>{farmer.name}</li>
        ))}
      </ul>
    </div>
  );
};

export default FarmersList;
