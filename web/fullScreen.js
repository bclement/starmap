function toEquitorial(lon, lat) {
    var ra = (-lon) / 15.0 + 12.0; 
    var decl = lat;
    return new OpenLayers.LonLat(ra, decl);
}
function toStellar(bounds) {
    var rval = new OpenLayers.Bounds();
    rval.extend(toEquitorial(bounds.left, bounds.bottom));
    rval.extend(toEquitorial(bounds.right, bounds.top));
    return rval;
}
var mousePositionCtrl = new OpenLayers.Control.MousePosition({
    prefix: "Equitorial coordinates: ",
    formatOutput: function(lonlat) {
        eq = toEquitorial(lonlat.lon, lonlat.lat);
        return "( " +eq.lon + ", " + eq.lat + " )";
    }
}
);
layer = new OpenLayers.Layer.WMS( "OpenLayers WMS", "/wms",
        {layers: 'stars'} );
layer.getURL = function (bounds) {
    bounds = this.adjustBounds(bounds);
    bounds = toStellar(bounds)

        var imageSize = this.getImageSize();
    var newParams = {};
    // WMS 1.3 introduced axis order
    var reverseAxisOrder = this.reverseAxisOrder();
    newParams.BBOX = this.encodeBBOX ?
        bounds.toBBOX(null, reverseAxisOrder) :
        bounds.toArray(reverseAxisOrder);
    newParams.WIDTH = imageSize.w;
    newParams.HEIGHT = imageSize.h;
    var requestString = this.getFullRequestString(newParams);
    return requestString;
};
var map = new OpenLayers.Map({
    div: "map",
    layers: [layer],
    /*
    controls: [
        new OpenLayers.Control.Navigation({
            dragPanOptions: {
                enableKinetic: true
            }
        }),
        new OpenLayers.Control.PanZoom(),
        new OpenLayers.Control.Attribution()
    ],*/
    center: [0, 0],
    zoom: 3
});

map.addControl(new OpenLayers.Control.LayerSwitcher());
map.addControl(mousePositionCtrl);
