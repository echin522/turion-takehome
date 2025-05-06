import { Center, Flex, Table, Text } from "@mantine/core";
import { useCallback, useEffect, useState } from "react";

const NO_UNIT = "N/A"

export interface RawPacket {
  ccsdsPrimaryHeader: {
    packetId: number;
    packetSeqCtrl: number;
    packetLength: number;
  };
  ccsdsSecondaryHeader: {
    timestamp: number;
    subsystemId: number;
  };
  telemetryPayload: {
    temperature: number;
    battery: number;
    altitude: number;
    signal: number;
  };
}

interface FieldValueRow {
  field: string;
  value: string | number;
  unit: string;
}

interface CurrentTelemetryTableProps {
  height?: number;
}

export const CurrentTelemetryTable = ({height = 400}: CurrentTelemetryTableProps) => {
  const [data, setData]       = useState<FieldValueRow[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError]     = useState<string | null>(null);

  const fetchCurrentTelemetry = useCallback(() => {
    fetch("/api/telemetry/current")
      .then(async (res) => {
        if (!res.ok) {
          const body = await res.text();
          throw new Error(`HTTP ${res.status}: ${body}`);
        }
        return res.json() as Promise<RawPacket>;
      })
      .then((pkt) => {
        setError(null);
        // build one row per field
        const rows: FieldValueRow[] = [
          { 
            field: "Packet ID",
            value: pkt.ccsdsPrimaryHeader.packetId,
            unit: NO_UNIT,
          },
          { 
            field: "Seq Ctrl",
            value: pkt.ccsdsPrimaryHeader.packetSeqCtrl,
            unit: NO_UNIT,
          },
          { 
            field: "Packet Length",
            value: pkt.ccsdsPrimaryHeader.packetLength,
            unit: "Bytes",
          },
          { 
            field: "Timestamp",
            value: new Date(pkt.ccsdsSecondaryHeader.timestamp * 1000).toISOString(),
            unit: NO_UNIT,
          },
          { 
            field: "Subsystem ID",
            value: pkt.ccsdsSecondaryHeader.subsystemId,
            unit: NO_UNIT,
          },
          { 
            field: "Temperature",
            value: pkt.telemetryPayload.temperature,
            unit: "C",
          },
          { 
            field: "Battery",
            value: pkt.telemetryPayload.battery,
            unit: "%",
          },
          { 
            field: "Altitude",
            value: pkt.telemetryPayload.altitude,
            unit: "km",
          },
          { 
            field: "Signal",
            value: pkt.telemetryPayload.signal,
            unit: "dB",
          },
        ];
        setData(rows);
      })
      .catch((err: any) => {
        console.error(err);
        setError(err.message || "Unknown error");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [])

  useEffect(() => {
    setLoading(true);
    fetchCurrentTelemetry();
    const iv = setInterval(fetchCurrentTelemetry, 1000);
    return () => clearInterval(iv);
  }, [fetchCurrentTelemetry]);

  // TODO: MRT is not playing nicely with styling likely because 2.0.0 uses v7 of
  // Mantine, but I need a feature from v8 of Mantine. I forget how I fixed this
  // Before So I will just create a table the old fashioned way
  // const columns: MRT_ColumnDef<FieldValueRow>[] = [
  //   { accessorKey: "field", header: "Field" },
  //   { accessorKey: "value", header: "Value" },
  //   { accessorKey: "unit",  header: "Unit" },
  // ];

  if (error) {
    return (
      <Center style={{ height: 200 }}>
        <Text c="red">Error: {error}</Text>
      </Center>
    );
  }

  const rows = data.map(({ field, value, unit }) => (
    <tr key={field}>
      <td>{field}</td>
      <td>{value}</td>
      <td>{unit}</td>
    </tr>
  ));

  return (
    <Flex direction="column" h={height} style={{ alignItems:"center"}}>
      <Text size="lg" w={500} style={{ textAlign: "center" }}>Current Values</Text>
      <Table
        striped
        highlightOnHover
        horizontalSpacing="md"
        verticalSpacing="sm"
        w={500}
      >
        <thead>
          <tr>
            <th style={{ textAlign: "start" }}>Field</th>
            <th style={{ textAlign: "start" }}>Value</th>
            <th style={{ textAlign: "start" }}>Unit</th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </Table>
    </Flex>
  );
}
