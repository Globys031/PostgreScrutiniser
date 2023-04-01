// Open/Close content. This function is necessary for smooth collapse transition
export const toggleCollapseContent = (content: HTMLElement | null) => {
  if (content) {
    const scrollHeight = content.scrollHeight;
    const maxHeight = content.style.maxHeight;

    maxHeight === "" || maxHeight === "0px"
      ? (content.style.maxHeight = scrollHeight + "px")
      : (content.style.maxHeight = "0");
  } else {
    console.error("content ref is undefined");
  }
};

// Only resize content if the collapsible is already open.
export const resizeContentIfOpen = (content: HTMLElement | null) => {
  if (content) {
    const scrollHeight = content.scrollHeight;
    const maxHeight = content.style.maxHeight;

    if (maxHeight === "" || maxHeight === "0px") {
      return;
    }
    content.style.maxHeight = scrollHeight + "px";
  } else {
    console.error("content ref is undefined");
  }
};
