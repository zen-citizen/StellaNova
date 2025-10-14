type Attribute = {
  name: string
  value: string
}

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
              className="text-secondary text-body-m font-medium pr-8 align-top text-left"
            >
              {attribute.name}
            </th>
            <td className="text-primary text-body-m font-semibold">
              {attribute.value}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

export default Table
