// components/TelemetryDashboard.tsx
import { Center, Flex, Group, Loader, Paper, Text } from '@mantine/core';
import { DateTimePicker, DateTimePickerProps, TimePickerProps } from '@mantine/dates';
import React, { useEffect, useState } from 'react';
import { AnomalyNotifications } from './AnomalyNotifications';
import { TelemetryChart, TelemetryRow } from './Chart';
import { CurrentTelemetryTable, RawPacket } from './Table';

const DefaultDateTimePickerProps: Partial<DateTimePickerProps> = {
  clearable: true,
  withSeconds: true,
  valueFormat: "DD MMM YYYY hh:mm",
}

const DefaultTimePickerProps: Partial <TimePickerProps> = {
  withDropdown: true,
  popoverProps: { 
    withinPortal: false,
    styles: {
      dropdown: {
        backgroundColor: "gray"
      }
    }
  },
  withSeconds: true,
  hoursStep: 1,
  minutesStep: 1,
  secondsStep: 1,
}

export const TelemetryDashboard: React.FC = () => {
  const now = new Date();
  const [start, setStart] = useState<Date>(new Date(now.getTime() - 3_600_000));
  const [end, setEnd]     = useState<Date>(now);

  const [data, setData]       = useState<TelemetryRow[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError]     = useState<string | null>(null);

  // fetch telemetry data
  useEffect(() => {
    let isMounted = true;
    const fetchHistoricalData = async () => {
      setLoading(true);
      setError(null);

      const qs = new URLSearchParams({
        start_time: start.toISOString(),
        end_time:   end.toISOString(),
      });
      
      try {
        const r = await fetch(`/api/telemetry?${qs.toString()}`);
        if (!r.ok) throw new Error(await r.text());
        const raw = (await r.json()) as Array<RawPacket>;
        const rows: TelemetryRow[] = raw.map((p) => ({
          temperature:   p.telemetryPayload.temperature,
          battery:       p.telemetryPayload.battery,
          altitude:      p.telemetryPayload.altitude,
          signal:        p.telemetryPayload.signal,
          time:          new Date(p.ccsdsSecondaryHeader.timestamp * 1000),
        }));
        setData(rows);
        
      } catch (err: any) {
        console.error(err);
        setError(err.message || 'Unknown error');
      } finally {
        setLoading(false);
      }
    }

    fetchHistoricalData();

    return () => {
      isMounted = false;
    }
  }, [start, end]);

  return (
    <Paper p="md" shadow="sm" bg="dark">
      <Group gap="md" mb="md" h="60px">
        <DateTimePicker
          {...DefaultDateTimePickerProps}
          timePickerProps={{...DefaultTimePickerProps}}
          label="Start Time"
          value={start}
          onChange={(d) => d && setStart(new Date(d))}
        />
        <DateTimePicker
          {...DefaultDateTimePickerProps}
          timePickerProps={{...DefaultTimePickerProps}}
          label="End Time"
          value={end}
          onChange={(d) => d && setEnd(new Date(d))}
        />
        {/* <Box style={{ height: "100%", position: "relative" }}>
          <Button 
            pos="absolute"
            bottom={1} 
            onClick={fetchHistoricalData}
            bg="gray"
          >
            Refresh
          </Button>
        </Box> */}
      </Group>

      {/* loading or error */}
      {error && (
        <Center my="md">
          <Text c="red">Error: {error}</Text>
        </Center>
      )}
      {loading ? (
        <Center style={{ height: 200 }}>
          <Loader />
        </Center>
      ) : (
        <Flex>
          <TelemetryChart data={data} />
          <CurrentTelemetryTable />
          <AnomalyNotifications />
        </Flex>
      )}
    </Paper>
  );
};
