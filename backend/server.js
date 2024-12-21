// /backend/server.js or /backend/index.js

const express = require('express');
const cors = require('cors');
const app = express();

// CORS configuration
app.use(cors({
  origin: 'http://localhost:3000',
  credentials: true
}));

app.use(express.json());

// Farmers routes
app.get('/farmers', async (req, res) => {
  try {
    // If you're using a database, replace this with actual database query
    const farmers = [
      { id: 1, name: 'John Doe' },
      { id: 2, name: 'Jane Smith' }
    ];
    res.json(farmers);
  } catch (error) {
    res.status(500).json({ message: 'Error fetching farmers', error: error.message });
  }
});

// Add more routes
app.post('/farmers', async (req, res) => {
  try {
    // Add logic to create a new farmer
    const newFarmer = req.body;
    // Save to database
    res.status(201).json(newFarmer);
  } catch (error) {
    res.status(500).json({ message: 'Error creating farmer', error: error.message });
  }
});

app.get('/farmers/:id', async (req, res) => {
  try {
    const { id } = req.params;
    // Add logic to fetch specific farmer
    res.json({ id, name: 'John Doe' });
  } catch (error) {
    res.status(500).json({ message: 'Error fetching farmer', error: error.message });
  }
});

app.put('/farmers/:id', async (req, res) => {
  try {
    const { id } = req.params;
    const updateData = req.body;
    // Add logic to update farmer
    res.json({ id, ...updateData });
  } catch (error) {
    res.status(500).json({ message: 'Error updating farmer', error: error.message });
  }
});

app.delete('/farmers/:id', async (req, res) => {
  try {
    const { id } = req.params;
    // Add logic to delete farmer
    res.json({ message: 'Farmer deleted successfully' });
  } catch (error) {
    res.status(500).json({ message: 'Error deleting farmer', error: error.message });
  }
});

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
