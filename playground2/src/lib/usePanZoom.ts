// Wheel zoom + drag pan for an SVG whose viewBox is controlled state. Shared by the world and explain graphs.
import { useCallback, useEffect, useRef, useState, type PointerEvent as ReactPointerEvent } from "react";

export interface View { x: number; y: number; w: number; h: number }

const DRAG_THRESHOLD = 4; // px before a press becomes a pan (below it, it is a click)

export function usePanZoom(home: View) {
  const [view, setView] = useState<View>(home);
  useEffect(() => setView(home), [home]);

  const svgRef = useRef<SVGSVGElement>(null);
  const drag = useRef<{ x: number; y: number; view: View; active: boolean } | null>(null);
  const dragged = useRef(false);

  // Wheel zoom around the cursor (needs a non-passive native listener to prevent page scroll).
  useEffect(() => {
    const el = svgRef.current;
    if (!el) return;
    const onWheel = (ev: WheelEvent) => {
      ev.preventDefault();
      const r = el.getBoundingClientRect();
      const f = Math.exp(ev.deltaY * 0.0015);
      setView((v) => {
        const scale = Math.min(Math.max(f, 0.5), 2);
        const w = Math.min(Math.max(v.w * scale, home.w / 8), home.w * 4);
        const k = w / v.w;
        const fx = (ev.clientX - r.left) / r.width, fy = (ev.clientY - r.top) / r.height;
        return { x: v.x + v.w * fx * (1 - k), y: v.y + v.h * fy * (1 - k), w, h: v.h * k };
      });
    };
    el.addEventListener("wheel", onWheel, { passive: false });
    return () => el.removeEventListener("wheel", onWheel);
  }, [home]);

  const reset = useCallback(() => setView(home), [home]);

  const handlers = {
    onDoubleClick: reset,
    onPointerDown: (e: ReactPointerEvent<SVGSVGElement>) => {
      dragged.current = false;
      drag.current = { x: e.clientX, y: e.clientY, view, active: false };
    },
    onPointerMove: (e: ReactPointerEvent<SVGSVGElement>) => {
      const d = drag.current;
      if (!d) return;
      if (!d.active) {
        if (Math.hypot(e.clientX - d.x, e.clientY - d.y) < DRAG_THRESHOLD) return;
        d.active = true;
        dragged.current = true;
        e.currentTarget.setPointerCapture(e.pointerId);
      }
      const r = e.currentTarget.getBoundingClientRect();
      // meet: the same scale applies to both axes
      const s = Math.max(d.view.w / r.width, d.view.h / r.height);
      setView({ ...d.view, x: d.view.x - (e.clientX - d.x) * s, y: d.view.y - (e.clientY - d.y) * s });
    },
    onPointerUp: () => { drag.current = null; },
  };

  return { svgRef, viewBox: `${view.x} ${view.y} ${view.w} ${view.h}`, reset, handlers, wasDragged: () => dragged.current };
}
