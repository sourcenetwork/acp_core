import { useCallback, useEffect, useRef } from "react";

export function useDebounce<T extends (...args: Parameters<T>) => void>(
  callback: T,
  delay: number
): T {
  const timeoutRef = useRef<number | null>(null);
  const functionRef = useRef<{ callback: T; args: Parameters<T> } | null>(null);

  const debounce = useCallback(
    (...args: Parameters<T>) => {
      functionRef.current = { callback, args };
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
      timeoutRef.current = setTimeout(() => callback(...args), delay);
    },
    [callback, delay]
  );

  useEffect(() => {
    // Cleanup the timeout when the component unmounts or delay changes
    return () => {
      // timeoutRef.current && clearTimeout(timeoutRef.current);
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, [delay]);

  return debounce as T;
}
