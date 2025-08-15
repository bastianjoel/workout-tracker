class WorkoutBreakdown extends HTMLElement {
  constructor() {
    super();

    this.activeItem = null;

    this.mapEl = document.getElementById(this.getAttribute("map-id"));
    this.chartEl = document.getElementById(this.getAttribute("workout-stats"));

    this.data = JSON.parse(
      document.getElementById(this.getAttribute("data-el")).textContent,
    );
    this.preferredUnits = JSON.parse(
      document.getElementById(this.getAttribute("preferred-units-el"))
        .textContent,
    );

    this.availableMetrics = {
      time: "",
      distance: this.preferredUnits.distance,
      duration: "",
      speed: this.preferredUnits.speed,
      elevation: this.preferredUnits.elevation,
      "heart-rate": this.preferredUnits.heartRate,
      cadence: this.preferredUnits.cadence,
      temperature: this.preferredUnits.temperature,
    };
  }

  connectedCallback() {
    this.querySelectorAll(`.breakdown-item`).forEach((item) => {
      item.addEventListener("mouseover", () => this.itemMouseOver(item));
      item.addEventListener("mouseout", this.itemMouseOut.bind(this));
      item.addEventListener("click", () => this.itemClick(item));
    });

    this.render();
  }

  itemClick(item) {
    if (this.activeItem === item) {
      this.setActiveItem(null);
    } else {
      this.setActiveItem(item);
    }
  }

  itemMouseOver(item) {
    if (!this.activeItem) {
      this.mapEl.setMarker(item);
    }
  }

  itemMouseOut() {
    if (!this.activeItem) {
      this.mapEl.clearMarker();
    }
  }

  setActiveItem(item) {
    if (this.activeItem) {
      this.activeItem.classList.remove(`active`);
    }

    this.activeItem = item;
    if (this.activeItem) {
      this.mapEl.scrollIntoView({ behavior: `smooth` });
      this.activeItem.classList.add(`active`);
      this.mapEl.setMarker(this.activeItem);
    } else {
      this.mapEl.clearMarker();
    }
  }

  render() {
    const header = this.renderHeader();
    this.innerHTML = `
      <table class="breakdown-table">
        <thead></thead>
        <tbody></tbody>
      </table>
    `;
    this.querySelector("thead").appendChild(header);
  }

  renderHeader() {
    const header = document.createElement("tr");
    header.classList.add("breakdown-header");

    console.log("Available Metrics:", this.availableMetrics);
    for (const metric of Object.keys(this.availableMetrics)) {
      if (this.data[metric] !== undefined) {
        const col = this.data[metric].Label;
        if (col) {
          header.appendChild(document.createElement("th")).textContent = col;
        }
      }
    }

    return header;
  }
}

customElements.define("workout-breakdown", WorkoutBreakdown);
