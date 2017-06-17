// Updated and customized version of RGBaster.js

// Helper functions.
const getContext = (width, height) => {
  const canvas = document.createElement('canvas');
  canvas.setAttribute('width', width);
  canvas.setAttribute('height', height);
  return canvas.getContext('2d');
};

const getImageData = (img, loaded) => {
  const imgObj = new Image();
  const imgSrc = img.src || img;

  // Can't set cross origin to be anonymous for data url's
  // https://github.com/mrdoob/three.js/issues/1305
  if (imgSrc.substring(0, 5) !== 'data:') imgObj.crossOrigin = 'Anonymous';

  imgObj.onload = function() {
    const context = getContext(imgObj.width, imgObj.height);
    context.drawImage(imgObj, 0, 0);

    const imageData = context.getImageData(0, 0, imgObj.width, imgObj.height);
    loaded && loaded(imageData.data);
  };

  imgObj.src = imgSrc;
};

const makeRGB = name => {
  return ['rgb(', name, ')'].join('');
};

const mapPalette = palette => {
  const arr = [];
  for (var prop in palette) {
    arr.push(frmtPobj(prop, palette[prop]));
  }
  arr.sort(function(a, b) {
    return b.count - a.count;
  });
  return arr;
};

const fitPalette = (arr, fitSize) => {
  if (arr.length > fitSize) {
    return arr.slice(0, fitSize);
  } else {
    for (var i = arr.length - 1; i < fitSize - 1; i++) {
      arr.push(frmtPobj('0,0,0', 0));
    }
    return arr;
  }
};

const frmtPobj = (a, b) => {
  return { name: makeRGB(a), count: b };
};

const MAX_WHITE = 690; // < (230, 230, 230);
const MAX_BLACK = 150; // > (50, 50, 50);
const MIN_DIFFERENCE = 20; // prevent r==g==b to eclude blacks and whites

const findBestColor = palette => {
  let best = palette[0];
  for (var i = 0; i < palette.length; i++) {
    const color = palette[i].name;
    let sum = 0;
    let difference = 0; // prevent
    let colors = color
      .substring(4, color.length - 1)
      .split(',')
      .map(v => parseInt(v));

    colors.forEach(v => {
      sum += v;
      difference += Math.abs(colors[0] - v);
    });
    if (sum > MAX_BLACK && sum < MAX_WHITE) {
      if (difference > MIN_DIFFERENCE) {
        return palette[i];
      }
      best = palette[i];
    }
  }
  return best;
};

// RGBaster Object
// ---------------
//
const PALETTESIZE = 10;

export default {
  colors: (img, opts) => {
    opts = opts || {};
    const exclude = opts.exclude || ['rgb(0,0,0)', 'rgb(255,255,255)'], // for example, to exclude white and black
      paletteSize = opts.paletteSize || PALETTESIZE;

    getImageData(img, function(data) {
      let colorCounts = {},
        rgbString = '',
        rgb = [],
        colors = {
          dominant: { name: '', count: 0 },
          palette: []
        };

      let i = 0;
      for (; i < data.length; i += 4) {
        rgb[0] = data[i];
        rgb[1] = data[i + 1];
        rgb[2] = data[i + 2];
        rgbString = rgb.join(',');

        // skip undefined data and transparent pixels
        if (rgb.indexOf(undefined) !== -1 || data[i + 3] === 0) {
          continue;
        }

        // Ignore those colors in the exclude list.
        if (exclude.indexOf(makeRGB(rgbString)) === -1) {
          if (rgbString in colorCounts) {
            colorCounts[rgbString] = colorCounts[rgbString] + 1;
          } else {
            colorCounts[rgbString] = 1;
          }
        }
      }

      if (opts.success) {
        const palette = fitPalette(mapPalette(colorCounts), paletteSize + 1);
        opts.success({
          dominant: palette[0].name,
          secondary: palette[1].name,
          best: findBestColor(palette).name,
          palette: palette
            .map(function(c) {
              return c.name;
            })
            .slice(1)
        });
      }
    });
  }
};
