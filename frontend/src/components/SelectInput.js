import React from "react";

function SelectInput({ label, options, value, onChange }) {
  return (
    <div className="mb-4">
      <label className="block mb-1">{label}</label>
      <select
        value={value}
        onChange={onChange}
        className="border border-gray-300 p-2 rounded w-full"
      >
        {options.map((opt) => (
          <option key={opt} value={opt}>
            {opt}
          </option>
        ))}
      </select>
    </div>
  );
}

export default SelectInput;
