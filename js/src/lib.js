import mjml2html from "mjml-browser";
import htmlMinifier from "html-minifier";
import jsBeautify from "js-beautify";

const defaultMinifyOptions = {
  caseSensitive: true,
  collapseWhitespace: true,
  minifyCSS: false,
  removeEmptyAttributes: true,
};

const defaultBeautifyOptions = {
  end_with_newline: true,
  indent_size: 2,
  max_preserve_newline: 0,
  preserve_newlines: false,
  wrap_attributes_indent_size: 2,
};

export function compile(input) {
  if (!input.mjml) {
    return {
      error: {
        message: "input is missing mjml property",
      },
    };
  }

  let options = {};

  if (input.options) {
    options = omit(input.options, "beautify", "minify", "minifyOptions");
  }

  let output;

  try {
    output = mjml2html(input.mjml, options);

    if (input.options && input.options.beautify) {
      if (input.options.beautifyOptions == null) {
        input.options.beautifyOptions = {};
      }

      output.html = jsBeautify.html(output.html, {
        ...defaultBeautifyOptions,
        ...input.options.beautifyOptions,
      });
    }

    if (input.options && input.options.minify) {
      if (input.options.minifyOptions == null) {
        input.options.minifyOptions = {};
      }

      // Convert array of regex strings to RegExp objects
      const regexArrays = [
        "customAttrAssign",
        "customAttrSurround",
        "ignoreCustomComments",
        "ignoreCustomFragments",
      ];

      regexArrays.forEach(function (key) {
        if (
          input.options.minifyOptions &&
          input.options.minifyOptions.hasOwnProperty(key)
        ) {
          input.options.minifyOptions[key] = parseStringRegExpArray(
            input.options.minifyOptions[key]
          );
        }
      });

      // Convert regex string to RegExp object
      if (
        input.options.minifyOptions &&
        input.options.minifyOptions.customAttrCollapse
      ) {
        input.options.minifyOptions.customAttrCollapse = parseRegExp(
          input.options.minifyOptions.customAttrCollapse
        );
      }

      output.html = htmlMinifier.minify(output.html, {
        ...defaultMinifyOptions,
        ...input.options.minifyOptions,
      });
    }
  } catch (err) {
    return {
      error: {
        message: err.message,
      },
    };
  }

  let result = {
    html: output.html,
  };

  if (output.errors.length > 0) {
    result.error = {
      message: "MJML compilation error",
      details: output.errors.map((err) => ({
        line: err.line,
        message: err.message,
        tagName: err.tagName,
      })),
    };
  }

  return result;
}

function omit(obj, ...props) {
  const result = { ...obj };

  props.forEach(function (prop) {
    delete result[prop];
  });

  return result;
}

function parseRegExp(value) {
  if (value) {
    return new RegExp(value.replace(/^\/(.*)\/$/, "$1"));
  }
}

function parseStringRegExpArray(value) {
  if (!Array.isArray(value)) {
    value = [value];
  }

  return value && value.map(parseRegExp);
}
