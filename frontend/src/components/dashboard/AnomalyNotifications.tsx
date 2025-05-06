import { notifications } from "@mantine/notifications";
import { useCallback, useEffect, useRef, useState } from "react";

interface Anomaly {
  field: number;  // Enum of field, not field name itself
  value: number;
  timestamp: number; // UNIX seconds
}

// TODO: This should have been generated with protobufs but protoc-gen-ts has not
// been kind to me
const AnomalyMap: Record<number, string> = {
	1: "TEMPERATURE",
  2: "BATTERY",
  3: "ALTITUDE",
  4: "SIGNAL",
}

export const AnomalyNotifications = () => {
  const [error, setError] = useState<string | null>(null);
  
  // keep track of which anomalies we've already notified
  // this is a memory leak waiting to happen
  const notifiedRef = useRef<Set<string>>(new Set());
  
  // Fetch a running window of 3 seconds for anomalies and display new ones to
  // the user. I recognize this is prone to race conditions - I would instead 
  // use an alerting mechanism but I want to get something out now
  const fetchAnomalies = useCallback(async () => {
    try {
      const currTime = new Date()
      const qs = new URLSearchParams({
        start_time: new Date(currTime.getTime() - 5_000).toISOString(),
        end_time:   currTime.toISOString(),
      });
      const r = await fetch(`/api/telemetry/anomalies?${qs.toString()}`);
      if (!r.ok) throw new Error(await r.text());
      const anomalies = (await r.json()) as Anomaly[];
      anomalies.forEach((an) => {
        const key = `${AnomalyMap[an.field]}@${an.timestamp}`;
        if (!notifiedRef.current.has(key)) {
          notifiedRef.current.add(key);
          notifications.show({
            title: `Anomaly: ${AnomalyMap[an.field]}`,
            message: `Value ${an.value} at ${new Date(an.timestamp * 1000).toLocaleTimeString()}`,
            color: 'red',
            autoClose: 3000,
          });
        }
      });
      setError(null);
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Unknown error');
    }
  }, [])

  useEffect(() => {
    const iv = setInterval(fetchAnomalies, 2000);
    return () => clearInterval(iv);
  }, [fetchAnomalies]);
  // TODO: Show an error notification
  return <></>
}