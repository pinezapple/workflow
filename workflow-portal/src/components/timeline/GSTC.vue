<template>
  <div>
    <div ref="gstc" class="gstc-wrapper" />
  </div>
</template>

<script>
import moment from "moment";
import GSTC from "gantt-schedule-timeline-calendar";
import { Plugin as TimelinePointer } from "gantt-schedule-timeline-calendar/dist/plugins/timeline-pointer.esm.min.js";
import { Plugin as Selection } from "gantt-schedule-timeline-calendar/dist/plugins/selection.esm.min.js";
import { Plugin as ItemResizing } from "gantt-schedule-timeline-calendar/dist/plugins/item-resizing.esm.min.js";
import { Plugin as ItemMovement } from "gantt-schedule-timeline-calendar/dist/plugins/item-movement.esm.min.js";
import "gantt-schedule-timeline-calendar/dist/style.css";
let gstc, state;
// helper functions
/**
 * @returns { import("gantt-schedule-timeline-calendar").Rows }
 */

function generateRows(tasks) {
  if (tasks.length === 0) return {};

  const rows = {};
  for (let i = 0; i < tasks.length; i++) {
    const id = GSTC.api.GSTCID((i + 1).toString());
    rows[id] = {
      id,
      label: `<a href='javascript:void(0);'  onclick="window.location.href = '/tasks/${tasks[i].id}'" target="_self" style="color:#0077c0; ">${tasks[i].name}</a>`,
      log: `<a href='javascript:void(0);'  onclick="window.location.href = '/log/${tasks[i].id}'" target="_self" style="color:#0077c0; ">Logs</a>`,
    };
  }
  return rows;
}
function generateItems(tasks) {
  if (tasks.length === 0) return {};

  const items = {};
  for (let i = 0; i < tasks.length; i++) {
    const id = GSTC.api.GSTCID((i + 1).toString());
    const rowId = GSTC.api.GSTCID((i + 1).toString());
    items[id] = {
      id,
      label: tasks[i].name,
      rowId,
      time: {
        start: GSTC.api.date(moment(tasks[i].started_at)).valueOf(),
        end: GSTC.api.date(moment(tasks[i].end_at)).valueOf(),
      },
    };
  }
  return items;
}

function generateChartTime(tasks) {
  const sortedTasks = tasks.sort((a, b) => a.started_at - b.started_at);
  const ntasks = sortedTasks.length;
  if (ntasks === 0) {
    return {
      from: moment().valueOf(),
      to: moment().valueOf(),
    };
  }
  return {
    from: GSTC.api
      .date(moment(sortedTasks[0].started_at).add(-5, "seconds"))
      .valueOf(),
    to: GSTC.api
      .date(moment(sortedTasks[ntasks - 1].end_at).add(5, "seconds"))
      .valueOf(),
    calculatedZoomMode: true,
  };
}

function generateCalendarLevels(tasks) {
  const sortedTasks = tasks.sort((a, b) => a.started_at - b.started_at);
  const ntasks = sortedTasks.length;
  const overDuration = moment(sortedTasks[ntasks - 1].end_at).diff(
    sortedTasks[0].started_at,
    "seconds"
  );

  const months = [
    {
      zoomTo: 100,
      period: "month",
      main: true,
      format({ timeStart }) {
        return timeStart.format("MM/YYYY");
      },
    },
  ];

  const days = [
    {
      zoomTo: 100,
      period: "day",
      main: true,
      format({ timeStart }) {
        return timeStart.format("DD/MM/YYYY");
      },
    },
  ];

  const fullHours = [
    {
      zoomTo: 100,
      period: "hour",
      main: true,
      format({ timeStart }) {
        return timeStart.format("DD/MM/YYYY HH");
      },
    },
  ];

  const hours = [
    {
      zoomTo: 100,
      period: "hour",
      main: true,
      format({ timeStart }) {
        return timeStart.format("HH");
      },
    },
  ];

  const fullMinutes = [
    {
      zoomTo: 100,
      period: "minute",
      main: true,
      format({ timeStart }) {
        return timeStart.format("DD/MM/YYYY HH:mm");
      },
    },
  ];

  const minutes = [
    {
      zoomTo: 100,
      period: "minute",
      main: true,
      format({ timeStart }) {
        return timeStart.format("mm");
      },
    },
  ];

  const seconds = [
    {
      zoomTo: 100,
      period: "second",
      main: true,
      format({ timeStart }) {
        return timeStart.format("HH:mm:ss");
      },
    },
  ];

  if (overDuration < 60) {
    return [fullMinutes, seconds];
  } else if (overDuration < 3600) {
    return [fullHours, minutes];
  } else if (overDuration < 86400) {
    return [days, hours];
  } else return [months, days];
}
// main component
export default {
  name: "GSTC",
  props: {
    items: {
      type: Array,
      required: true,
      default: () => [],
    },
    showIndex: {
      type: Boolean,
      required: false,
      default: true,
    },
  },
  updated() {
    let columnsData = {};
    if (this.showIndex) {
      columnsData = {
        [GSTC.api.GSTCID("id")]: {
          id: GSTC.api.GSTCID("id"),
          width: 60,
          data: ({ row }) => GSTC.api.sourceID(row.id),
          header: { content: "ID" },
        },
        [GSTC.api.GSTCID("label")]: {
          id: GSTC.api.GSTCID("label"),
          width: 200,
          data: "label",
          header: { content: "Label" },
          isHTML: true,
        },
        [GSTC.api.GSTCID("log")]: {
          id: GSTC.api.GSTCID("log"),
          width: 70,
          data: "log",
          header: { content: "Log" },
          isHTML: true,
        },
      };
    } else {
      columnsData = {
        [GSTC.api.GSTCID("label")]: {
          id: GSTC.api.GSTCID("label"),
          width: 200,
          data: "label",
          header: { content: "Label" },
          isHTML: true,
        },
        [GSTC.api.GSTCID("log")]: {
          id: GSTC.api.GSTCID("log"),
          width: 70,
          data: "log",
          header: { content: "Log" },
          isHTML: true,
        },
      };
    }

    const config = {
      licenseKey:
        "====BEGIN LICENSE KEY====\nQaleTkiwAfi3UNBzBM76AGf/TrLhv93ldB4Pw67b0/Q9vVIsOydCWb/xo0ucOINUaOrud5Pa/iJZOAl/SFTCT1cu7qs01bDmlyMv93t2E+o8E2DtNngPKUIfVJn5vMBhJPbSpAs5A131l/QTtOMpXLAalCv7WidzRe9QesDzcOnqEPBUnHrZfXTUA7MF3avnO/Wx/nS3A+ZkFeIlhVS3Az+3oTKaX297uFOO4yjaYEvGRA1T4eXj1PHa4o2FSEe5+7Zg2cEfNgVgmURQrVt+VIsrw5KdAq6wpPJlLW1iGXBeOYJ4UwRaMlvheszIEkm0mgEKjc0zkRwZ2tXkWGFroA==||U2FsdGVkX18KGxxK0jt+gouSOI+jw1QrvwjCDsNGdo6uCyiW1XN4mLlCBFIz7/fa2P0O02tpu1wz7FjNymRWcAmhJLqb0usLJf5j0TepEJU=\ns1svqejqckhY0euXCvTecAKVcea0qUVNA5WF+BK5rGkugxiQFLucBNlFTVZ74+9ve9V8tbOXycVd4V9zuangZ41uQnLcsgC0RjTifDn7vTwbr8ehJDPIBQUimK7eJjAGxP0GjeDyZs0eL7S+Uk51UhEIa2ixvyC02Z3TEtcUHUqZTzdidruswvyxSKTCJAA+WJZ+7N5Bn27E0jj6gHxsjsFmRzUqQm4t6bY4FeBtG1YDTKhRl6z0ZT69eeqkwKnKYCDZc/IEx82cuNqXmrZBy/QKfDqtructkh0SErAuYr4Hyubjv38XY2cZnUx+kQHRArPpddzzxs0VoK3iIuSz2g==\n====END LICENSE KEY====",
      plugins: [TimelinePointer(), Selection(), ItemResizing(), ItemMovement()],
      innerHeight: 300,
      list: {
        expander: {
          straight: false,
        },
        columns: {
          data: columnsData,
        },
        rows: generateRows(this.items),
      },
      chart: {
        items: generateItems(this.items),
        time: generateChartTime(this.items),
        calendarLevels: generateCalendarLevels(this.items),
      },
      actions: {
        "chart-timeline-items-row-item": [this.itemClickAction],
      },
    };
    state = GSTC.api.stateFromConfig(config);
    gstc = GSTC({
      element: this.$refs.gstc,
      state,
    });
  },
  beforeUnmount() {
    if (gstc) gstc.destroy();
  },
  methods: {
    updateFirstRow() {
      state.update(`config.list.rows.${GSTC.api.GSTCID("0")}`, (row) => {
        row.label = "Changed dynamically";
        return row;
      });
    },
    changeZoomLevel() {
      state.update("config.chart.time.zoom", 17);
    },
    itemClickAction(element, data) {
      function onClick() {
        this.$router.push({
          name: "item_page_name",
          params: { id: data.item.id },
        });
      }
      element.addEventListener("click", onClick);
    },
  },
};
</script>
<style scoped>
.gstc-component {
  margin: 0;
  padding: 0;
}
.toolbox {
  padding: 10px;
}
</style>
