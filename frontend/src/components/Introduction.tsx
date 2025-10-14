function Introduction() {
  return (
    <>
      <section className="prose">
        <p>
          If you're a Bengaluru resident, you can use Civic Compass to identify
          the BBMP, BDA, Revenue, BESCOM, BWSSB offices, and Police stations for
          your area.
        </p>

        <p>
          <strong>Enter the exact address you need information for.</strong>
        </p>

        <p>This tool is only for Bengaluru at this time.</p>
      </section>

      <section className="prose mt-6">
        <h2 className="text-primary text-heading-xs">Data sources</h2>
        <p>
          We pull information from Government records. While we strive for
          accuracy, these sources can sometimes be incomplete or outdated.
        </p>
        <ul className="prose">
          <li>
            <a className="text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
              OpenCity
            </a>
          </li>
          <li>
            <a className="text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
              Karnataka-GIS
            </a>
          </li>
          <li>
            <a className="text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer">
              OpenStreetMap
            </a>
          </li>
        </ul>
      </section>
    </>
  )
}

export default Introduction
