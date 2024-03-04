import { Color } from '@tiptap/extension-color'
import { Editor } from '@tiptap/core'
import ListItem from '@tiptap/extension-list-item'
import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight'
import Image from '@tiptap/extension-image'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import Youtube from '@tiptap/extension-youtube'

const lowlight = require('./unexported/lowlight.js').lowlight;
const hashtag = require('./unexported/hashtag.js').HashTag;

const editorProps = require('./unexported/editor_props.js').EditorProps;

export const editor = new Editor({
   element: document.querySelector('#code-editor'),
   extensions: [
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
      CodeBlockLowlight.extend({
         addKeyboardShortcuts() {
            return {
               Tab: () => {
                  const { state, dispatch } = this.editor;
                  const { selection } = state;
                  const node = state.doc.nodeAt(selection.$from.pos);

                  if (node.type.name === 'codeBlock') {
                     const text = node.textContent;
                     const lines = text.split('\n');
                     const currentLine = lines[selection.$from.pos - node.pos];
                     const indent = currentLine.match(/^\s*/)[0];
                     const newIndent = indent + '\t';

                     const transaction = state.tr.insertText(newIndent);
                     dispatch(transaction);

                     return true;
                  }

                  return false;
               }
            };
         }
      }).configure({ 
         lowlight,
         HTMLAttributes: {
            class: 'mockup-code rounded-badge',
         }
      }),
      Image.configure({
         inline: true,
         HTMLAttributes: {
            class: 'max-h-96 mx-auto'
         }
      }), 
      Link.configure({
        openOnClick: true,
        autolink: true,
        linkOnPaste: true,
      }),
      Placeholder.configure({
         considerAnyAsEmpty: true,
         placeholder: 'Write your thought here...',
      }),
      Youtube.configure({
         inline: false,
         width: 512,
         height: 288,
         modestBranding: 'true',
         progressBarColor: 'white',
         HTMLAttributes: {
            class: 'mx-auto'
         }
      }),
      hashtag,
   ],
   editorProps: editorProps,
   onCreate: () => {
      const encodedContent = document.getElementById('code-editor').children[0].innerHTML
      document.getElementById('content').value = encodedContent
   },

   onUpdate: ({ editor }) => {
      const encodedContent = document.getElementById('code-editor').children[0].innerHTML
      console.log(encodedContent)
      document.getElementById('content').value = encodedContent

      const tags = editor.getJSON().content
         .filter((node) => node.type === 'paragraph')
         .filter((node) => node.content !== undefined)
         .filter((node) => node.content.length > 0)
         .flatMap((node) => 
            node
            .content
            .filter((child) => child !== undefined)
            .filter((child) => child.type === 'mention')
            .map((child) => child.attrs.id)
         )      

      const uqTags = [...new Set(tags)]
      document.getElementById('tags').value = uqTags.join(',')
   },
   onTransaction: (transaction) => {
      // const {state} = window.editor
      // const {selection} = state
      // const node = state.doc.nodeAt(selection.$from.pos)
      //
      // if (node && node.type.name === 'paragraph' || node.type.name === 'text'){
      //    const rect = editor.view.dom.getBoundingClientRect()
      //    const top = rect.top + window.pageYOffset
      //    const bottom = rect.bottom + window.pageYOffset
      //    const middle = (top + bottom) / 2
      //
      //    const range = document.createRange()
      //    range.setStart(node.dom, selection.$from.offset)
      //    range.collapse(true)
      //
      //    const caretRect = range.getClientRects()[0]
      //    const caretTop = caretRect.top + window.pageYOffset
      //    const caretBottom = caretRect.bottom + window.pageYOffset
      //
      //    if (caretTop < middle && middle + 50 < caretBottom) {
      //       node.dom.scrollIntoView({block: 'nearest', behavior: 'smooth', inline: 'nearest'})
      //    }
      // }
   },
   content: window.content,
   
})

window.editor = editor

window.editor.view.dom.addEventListener('paste', function(event) {
   event.preventDefault()

   const clipboardData = event.clipboardData || window.clipboardData
   const pastedContent = clipboardData.getData('text/html')

   const parsedContent = new DOMParser().parseFromString(pastedContent, 'text/html')

   editor.chain().focus().setContent(parsedContent.body.innerHTML).run()
})



window.allowDrop = function (ev) {
   ev.preventDefault()
}

window.upload = function(ev, url) {
   ev.preventDefault()

   if(ev.dataTransfer.files.length === 0) {
      return
   }

   file = ev.dataTransfer.files[0]

   Swal.showLoading()

   const formData = new FormData()
   formData.append("file", file)

   fetch(url, {
      method: "POST",
      body: formData,
      contentType: "multipart/form-data"
   }).then((response) => {
      return response.text()
   }).then(window.afterUpload);
}

window.afterUpload = function(rawData) {
   data = JSON.parse(rawData)
   editor.chain().focus().setImage({ src: data.data.url }).run()
   Swal.close()
}

