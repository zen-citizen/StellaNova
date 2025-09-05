import Map, { Marker } from "@vis.gl/react-maplibre";
import "maplibre-gl/dist/maplibre-gl.css";

// Bangalore coordinates
const bangaloreCoordinates = [12.9716, 77.5946]; // [latitude, longitude]

export const BlrMap = () => {
  return (
    <Map
      initialViewState={{
        longitude: bangaloreCoordinates[1],
        latitude: bangaloreCoordinates[0],
        zoom: 12,
      }}
      mapStyle="https://tiles.openfreemap.org/styles/liberty"
    >
      <Marker
        longitude={bangaloreCoordinates[1]}
        latitude={bangaloreCoordinates[0]}
      />
    </Map>
  );
};
