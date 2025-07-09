const axios = require("axios");
const { API_URL, API_HEADERS } = require("../../constants/index");
const errorHandler = require("../../utils/apiErrorHandler");

let rpcNode = null;
let initialized = false;

const initializeHandler = async () => {
  if (initialized) return;
  initialized = true;
  try {
    const response = await axios.get(API_URL, {
      headers: API_HEADERS,
    });
    rpcNode = response.data;
  } catch (error) {
    // Comment out problematic error handler that's causing server startup to fail
    // errorHandler(error.response?.data || error.message);
    console.log("External API call failed:", error.message);
  }
};

// Comment out the initialization that's causing startup issues
// Call the initialization
// initializeHandler();

// Export a higher-order function that wraps the module exports
const departmentModuleHandler = (moduleFactory) => {
  if (!initialized) {
    // initializeHandler();
  }
  return moduleFactory();
};

module.exports = { departmentModuleHandler };

