import type { Attribute } from "../types"

type TableProps = {
  attributes: Attribute[]
}

function Table(props: TableProps) {
  const { attributes } = props
  return (
    <table className="w-full border-separate border-spacing-y-4">
      <tbody>
        {attributes.map((attribute) => (
          <tr key={attribute.name}>
            <th
              scope="row"
              className="text-secondary text-body-m font-medium pr-8 align-top text-left w-1/2"
            >
              {attribute.name}
            </th>
            <td className="text-primary text-body-m flex flex-col gap-2">
              <span className="font-semibold">{attribute.value}</span>
              {attribute.address && (
                <>
                  <span>{attribute.address.text}</span>
                  <a
                    href={attribute.address.link}
                    target="_blank"
                    className="text-secondary underline underline-offset-3 hover:opacity-75 transition-opacity ease-in duration-100 cursor-pointer"
                  >
                    Google Maps
                  </a>
                </>
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export default Table
