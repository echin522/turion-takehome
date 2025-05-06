import React from "react";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export interface TelemetryRow {
  temperature: number;
  battery: number;
  altitude: number;
  signal: number;
  time: Date;
}

interface TelemetryChartProps {
  data: TelemetryRow[];
  height?: number;
}

export const TelemetryChart: React.FC<TelemetryChartProps> = ({
  data,
  height = 400,
}) => {
  const chartData = data.map((row)=> ({
    ...row,
    time: row.time.getTime(),
  }))
  return (
    <ResponsiveContainer width="100%" height={height}>
      <LineChart data={chartData} margin={{ top: 20, right: 50, bottom: 20, left: 20 }}>
        <CartesianGrid strokeDasharray="3 3" />

        {/* X axis: time in human‐readable form */}
        <XAxis
          dataKey="time"
          type="number"
          domain={['dataMin', 'dataMax']}
          tickFormatter={(ts) => new Date(ts).toLocaleTimeString()}
        />

        {/* Left‐side Y axis for temperature & altitude */}
        <YAxis
          yAxisId="left"
          label={{ value: 'Temp (°C) / Alt (m)', angle: -90, position: 'insideLeft' }}
        />

        {/* Right‐side Y axis for battery & signal */}
        <YAxis
          yAxisId="right"
          orientation="right"
          label={{ value: 'Battery (%) / Signal', angle: 90, position: 'insideRight' }}
        />

        <Tooltip
          labelFormatter={(ts) => new Date(ts).toLocaleString()}
        />
        <Legend verticalAlign="top" />

        {/* Temperature line */}
        <Line
          yAxisId="left"
          type="monotone"
          dataKey="temperature"
          name="Temperature (°C)"
          dot={false}
          stroke="#ff7300"
        />

        {/* Altitude line */}
        <Line
          yAxisId="left"
          type="monotone"
          dataKey="altitude"
          name="Altitude (m)"
          dot={false}
          stroke="#387908"
        />

        {/* Battery line */}
        <Line
          yAxisId="right"
          type="monotone"
          dataKey="battery"
          name="Battery (%)"
          dot={false}
          stroke="#8884d8"
        />

        {/* Signal line */}
        <Line
          yAxisId="right"
          type="monotone"
          dataKey="signal"
          name="Signal Strength"
          dot={false}
          stroke="#82ca9d"
        />
      </LineChart>
    </ResponsiveContainer>
  );
};
