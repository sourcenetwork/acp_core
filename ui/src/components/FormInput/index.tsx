import { Control, FieldValues, Path } from "react-hook-form";
import { FormControl, FormField, FormItem, FormMessage } from "../ui/form";
import { Input } from "../ui/input";

export const FormInputField = <T extends FieldValues>({ name, placeholder, control }: { name: Path<T>, placeholder: string, control: Control<T> }) => (
    <FormField
        control={control}
        name={name}
        render={({ field }) => (
            <FormItem>
                <FormControl>
                    <Input placeholder={placeholder} {...field} />
                </FormControl>
                <FormMessage />
            </FormItem>
        )}
    />
);
