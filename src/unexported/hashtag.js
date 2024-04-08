import Mention from '@tiptap/extension-mention'
import tippy from 'tippy.js'

export const HashTag = Mention.configure({
   renderHTML({ options, node}) {
      return [
         "span",
         { 
            class : "font-semibold text-primary/150 suggestion underline decoration-primary decoration-2",
            "data-type": "mention",
            "data-id": node.attrs.id,
         },
         `${options.suggestion.char}${node.attrs.label ?? node.attrs.id}`,
      ];
   },
   suggestion: {
      char: "#",
      render: () => {
         let popup
         let suggest

         let selectedIndex = 0

         return {
            onStart: (props) => {
               popup = tippy(document.body, {
                  getReferenceClientRect: props.clientRect,
                  appendTo: () => document.body,
                  content: `
                  <div class="dropdown dropdown-open">
                     <ul class="p-4 shadow menu dropdown-content z-[1] bg-base-200 rounded-box w-64">
                        <li class="text-lg text-base-content">
                           Start Write to Search...
                        </li>
                     </ul>
                  </div>
                  `,
                  allowHTML: true,
                  showOnCreate: true,
                  trigger: 'manual',
                  placement: 'bottom-start',
               })
            },
            onUpdate: (props) => {
               content = extractCharacterAfterHashtag(editor.state
                  .selection
                  .$head
                  .parent
                  .textContent
               )
            
               debounceUpdate(content, popup, function (headers) {
                  json = JSON.parse(headers.get('X-Tags'))
                  props.items = json
                  suggest = props
                  selectedIndex = 0
               })
            },
            onExit() {
               popup.destroy()
            },
            onKeyDown: (props) => {
               const { event } = props


               if (event.key === 'ArrowDown') {
                  event.preventDefault()
                  selectedIndex = refreshSuggestion(
                     selectedIndex === suggest.items.length-1 ? 0 : selectedIndex + 1
                  )
               }

               if (event.key === 'ArrowUp') {
                  event.preventDefault()
                  selectedIndex = refreshSuggestion(
                     selectedIndex === 0 ? suggest.items.length-1 : selectedIndex - 1,
                  )
               }

               if (event.key == 'Enter') {
                  event.preventDefault()
                  const item = suggest.items[selectedIndex]
                  suggest.command({ id:item })
               }

            },
         }
      }
   }
});

function refreshSuggestion(index) {
   getElementsByIdPrefix("tag-suggest-")
      .forEach((elem) => { elem.classList.remove("active") })

   document
      .querySelector(`#tag-suggest-${index}`)
      .classList
      .add("active")

   return index
}

function extractCharacterAfterHashtag(text) {
  const hashtagIndex = text.indexOf("#");

   if (hashtagIndex !== -1) {
      const whitespaceIndex = text.indexOf(" ", hashtagIndex + 1);

      if (whitespaceIndex !== -1) {
         return text.substring(hashtagIndex + 1, whitespaceIndex); 
      } else {
         return text.substring(hashtagIndex + 1, text.length);
      }
   }

  return null; 
}

function getElementsByIdPrefix(prefix) {
  const elements = document.querySelectorAll(`[id^="${prefix}"]`);
  return elements;
}

function debounce(callback, delay) {
   let timerID;
   return (...args) => {
      clearTimeout(timerID);
      timerID = setTimeout(() => {
         callback(...args);
      }, delay);
   }
}

const debounceUpdate = debounce((content, popup, onResponseHeader) => {
   fetch(`/dashboard/tags?tag=${content}`, {
      method: 'GET',
   }).then((response) => {
      onResponseHeader(response.headers)
      return response.text()
   }).then((data) => {
      if (popup.state.isDestroyed) return
      popup.setProps({
         content: data,
      })
   })
}, 500)

