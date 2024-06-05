interface StringObject {
	[key: string]: string;
}

export function addRemoveClassToElement(element: HTMLElement, addClass: string, removeClass: string) {
	element.classList.add(addClass);
	element.classList.remove(removeClass);
}

export function encodeAsFormData(obj: StringObject) {
	return new URLSearchParams(Object.entries(obj)).toString();
}
