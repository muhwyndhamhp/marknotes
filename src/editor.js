import { Color } from "@tiptap/extension-color";
import { Editor } from "@tiptap/core";
import ListItem from "@tiptap/extension-list-item";
import StarterKit from "@tiptap/starter-kit";
import TextStyle from "@tiptap/extension-text-style";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import Image from "@tiptap/extension-image";
import Link from "@tiptap/extension-link";
import Placeholder from "@tiptap/extension-placeholder";
import { Markdown } from "tiptap-markdown";
import BubbleMenu from "@tiptap/extension-bubble-menu";
import { PluginKey } from "@tiptap/pm/state";

const lowlight = require("./unexported/lowlight.js").lowlight;
const hashtag = require("./unexported/hashtag.js").HashTag;

const editorProps = require("./unexported/editor_props.js").EditorProps;

const yt = require("./unexported/youtube.js").Youtube;

export const editor = new Editor({
  element: document.querySelector("#code-editor"),
  extensions: [
    BubbleMenu.configure({
      shouldShow: ({ editor, from, to }) => {
        return from === to && editor.isActive("paragraph") && editor.isFocused;
      },
      pluginKey: new PluginKey("bubbleMenuOne"),
      tippyOptions: {
        placement: "bottom",
        delay: 1000,
      },
      element: document.querySelector(".quick-menu"),
    }),
    Markdown,
    Color.configure({ types: [TextStyle.name, ListItem.name] }),
    TextStyle.configure({ types: [ListItem.name] }),
    StarterKit.configure({
      bulletList: {
        keepMarks: true,
        keepAttributes: false,
      },
      orderedList: {
        keepMarks: true,
        keepAttributes: false,
      },
    }),
    CodeBlockLowlight.extend({}).configure({
      lowlight,
      HTMLAttributes: {
        class: "mockup-code rounded-badge",
      },
    }),
    Image.configure({
      inline: true,
      HTMLAttributes: {
        class: "max-h-96 mx-auto",
        loading: "lazy",
      },
    }),
    Link.configure({
      openOnClick: true,
      autolink: true,
      linkOnPaste: true,
    }),
    Placeholder.configure({
      considerAnyAsEmpty: true,
      placeholder: "Write your thought here...",
    }),
    yt.configure({
      inline: false,
      nocookie: true,
      modestBranding: "true",
      progressBarColor: "white",
      HTMLAttributes: {
        class:
          "mx-auto max-w-[320px] md:max-w-[512px] max-h-[240px] md:max-h-[288px]",
      },
    }),
    hashtag,
  ],
  editorProps: editorProps,
  onCreate: () => {
    const encodedContent =
      document.getElementById("code-editor").children[0].innerHTML;
    document.getElementById("content").value = encodedContent;
  },

  onUpdate: ({ editor }) => {
    const encodedContent =
      document.getElementById("code-editor").children[0].innerHTML;
    document.getElementById("content").value = encodedContent;

    const tags = editor
      .getJSON()
      .content.filter((node) => node.type === "paragraph")
      .filter((node) => node.content !== undefined)
      .filter((node) => node.content.length > 0)
      .flatMap((node) =>
        node.content
          .filter((child) => child !== undefined)
          .filter((child) => child.type === "mention")
          .map((child) => child.attrs.id)
      );

    const uqTags = [...new Set(tags)];
    document.getElementById("tags").value = uqTags.join(",");

    debouncePush();
  },
  content: window.content,
});

/**
 * @param {String} link
 * @param {String} text
 */
window.appendLinkedText = function (link, text) {
  editor.commands.insertContent(text);

  const { doc } = editor.state;
  const startNode = doc.content.size - 1 - text.length;
  const textNode = doc.textBetween(startNode, doc.content.size);

  if (textNode === text) {
    editor.commands.setTextSelection({
      from: startNode,
      to: doc.content.size,
    });
  }

  editor.commands.setLink({ link });
};

const debouncePush = debounce(() => {
  editor.commands.insertContent();
  st = document.getElementById("status").value;

  if (st && st === "published") {
    document.getElementById("publish-button").click();
  } else {
    if (document.getElementById("blog-title").value === "") {
      return;
    }

    document.getElementById("draft-button").click();
  }
}, 2_000); // 20 second

function debounce(callback, delay) {
  let timerID;
  return (...args) => {
    clearTimeout(timerID);
    timerID = setTimeout(() => {
      callback(...args);
    }, delay);
  };
}

window.editor = editor;

window.editor.view.dom.addEventListener("paste", async (event) => {
  const clipboardData = event.clipboardData || window.clipboardData;
  const pastedContent = clipboardData.getData("text/html");
  // const imageData = clipboardData.getData('image')

  if (
    pastedContent !== undefined &&
    pastedContent !== "" &&
    window.editor.isEmpty
  ) {
    event.preventDefault();
    const parsedContent = new DOMParser().parseFromString(
      pastedContent,
      "text/html"
    );

    editor.chain().focus().setContent(parsedContent.body.innerHTML).run();
    return;
  }

  // if (imageData !== undefined) {
  //    event.preventDefault()
  //    const blob = await new Promise((resolve) => {
  //       const reader = new FileReader();
  //       reader.onload = (event) => {
  //          resolve(event.target.result);
  //       };
  //       reader.readAsDataURL(imageData.getAsFile());
  //    });
  //
  //    // Create an URL object from the Blob
  //    const url = URL.createObjectURL(blob);
  //
  //    // Set the URL as the source of the image element
  //    console.log(url)
  // }
});
